package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/cloudbase/garm/params"
)

var doneChan = make(chan struct{})

func init() {
	close(doneChan)
}

func NewRunnerCommand(ctx context.Context, cmdParams []string, workdir string, forgeType params.EndpointType, st RunnerStateManager) (Worker, error) {
	runnerCfg := filepath.Join(workdir, ".runner")
	runCfg, err := NewRunnerConfig(runnerCfg, forgeType)
	if err != nil {
		return nil, fmt.Errorf("failed to read runner config: %w", err)
	}

	if len(cmdParams) < 1 {
		return nil, fmt.Errorf("cmdParams is empty")
	}

	mode, err := os.Stat(workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to access workdir: %w", err)
	}
	if !mode.IsDir() {
		return nil, fmt.Errorf("workdir %s is not a folder", workdir)
	}

	if err := os.Chdir(workdir); err != nil {
		return nil, fmt.Errorf("failed to chdir to %s: %w", workdir, err)
	}

	if st == nil {
		return nil, fmt.Errorf("invalid state manager")
	}
	command := exec.Command(cmdParams[0], cmdParams[1:]...)
	ret := &runnerCmd{
		ctx:       ctx,
		forgeType: forgeType,
		workdir:   workdir,
		cmd:       command,
		runnerCfg: runCfg,
		done:      doneChan,
		errChan:   make(chan error, 1),
		st:        st,
	}

	return ret, nil
}

type runnerCmd struct {
	ctx       context.Context
	forgeType params.EndpointType
	workdir   string
	runnerCfg RunnerConfig
	st        RunnerStateManager

	done    chan struct{}
	running bool
	mux     sync.Mutex

	cmd     *exec.Cmd
	cmdErr  error
	errChan chan error
}

func (r *runnerCmd) Wait() <-chan error {
	return r.errChan
}

func (r *runnerCmd) Start() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.running {
		return nil
	}

	r.done = make(chan struct{})
	r.running = true

	go r.loop()
	go r.executeCommand()
	return nil
}

func (r *runnerCmd) Stop() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if !r.running {
		return nil
	}

	close(r.done)
	r.running = false

	// Kill the command if it's running
	if r.cmd != nil && r.cmd.Process != nil {
		if err := r.cmd.Process.Kill(); err != nil {
			slog.ErrorContext(r.ctx, "failed to kill process", "error", err)
		}
	}

	return nil
}

func (r *runnerCmd) executeCommand() {
	var err error
	var stdout, stderr io.ReadCloser
	var jobStarted bool
	var jobMux sync.Mutex

	defer func() {
		r.cmdErr = err
		if stopErr := r.Stop(); stopErr != nil {
			slog.ErrorContext(r.ctx, "failed to stop runner", "error", stopErr)
		}
		r.st.SetRunnerStarted(false)

		jobMux.Lock()
		started := jobStarted
		jobMux.Unlock()

		if started {
			// Job was started and the runner exited. This means that the job reached a conclusion
			// and we need to remove the runner.
			slog.InfoContext(r.ctx, "runner has finished th job")
			r.st.SetJobFinished()
		}
	}()
	// Create pipes for stdout and stderr
	stdout, err = r.cmd.StdoutPipe()
	if err != nil {
		slog.ErrorContext(r.ctx, "failed to create stdout pipe", "error", err)
		return
	}

	stderr, err = r.cmd.StderrPipe()
	if err != nil {
		slog.ErrorContext(r.ctx, "failed to create stderr pipe", "error", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			slog.InfoContext(r.ctx, string(line))
			if r.isJobStartedLine(line) {
				slog.InfoContext(r.ctx, "runner is active")
				r.st.SetJobStarted()
				jobMux.Lock()
				jobStarted = true
				jobMux.Unlock()
				continue
			}

			if r.isRunnerStartedLine(line) {
				slog.InfoContext(r.ctx, "runner is online and idle")
				r.st.SetRunnerStarted(true)
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			slog.ErrorContext(r.ctx, "error reading stdout", "error", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Bytes()
			slog.InfoContext(r.ctx, string(line))
			if r.isJobStartedLine(line) {
				slog.InfoContext(r.ctx, "runner is active")
				r.st.SetJobStarted()
				jobMux.Lock()
				jobStarted = true
				jobMux.Unlock()
				continue
			}

			if r.isRunnerStartedLine(line) {
				slog.InfoContext(r.ctx, "runner is online and idle")
				r.st.SetRunnerStarted(true)
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			slog.ErrorContext(r.ctx, "error reading stderr", "error", err)
		}
	}()

	err = r.cmd.Start()
	if err != nil {
		slog.ErrorContext(r.ctx, "failed to start command", "error", err)
		return
	}

	err = r.cmd.Wait()
	if err != nil {
		slog.ErrorContext(r.ctx, "command failed", "error", err)
		return
	}
}

func (r *runnerCmd) loop() {
	defer func() {
		r.mux.Lock()
		defer r.mux.Unlock()
		r.errChan <- r.cmdErr
	}()

	for {
		select {
		case <-r.done:
			return
		case <-r.ctx.Done():
			return
		}
	}
}
