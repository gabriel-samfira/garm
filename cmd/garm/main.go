// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/cloudbase/garm-provider-common/util"
	"github.com/cloudbase/garm/apiserver/controllers"
	"github.com/cloudbase/garm/apiserver/routers"
	"github.com/cloudbase/garm/auth"
	"github.com/cloudbase/garm/config"
	"github.com/cloudbase/garm/database"
	"github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/metrics"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner" //nolint:typecheck
	runnerMetrics "github.com/cloudbase/garm/runner/metrics"
	garmUtil "github.com/cloudbase/garm/util"
	"github.com/cloudbase/garm/util/appdefaults"
	"github.com/cloudbase/garm/websocket"
)

var (
	conf    = flag.String("config", appdefaults.DefaultConfigFilePath, "garm config file")
	version = flag.Bool("version", false, "prints version")
)

var Version string

var signals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
}

func maybeInitController(db common.Store) error {
	if _, err := db.ControllerInfo(); err == nil {
		return nil
	}

	if _, err := db.InitController(); err != nil {
		return errors.Wrap(err, "initializing controller")
	}

	return nil
}

func setupLogging(ctx context.Context, logCfg config.Logging, hub *websocket.Hub) {
	logWriter, err := util.GetLoggingWriter(logCfg.LogFile)
	if err != nil {
		log.Fatalf("fetching log writer: %+v", err)
	}

	// rotate log file on SIGHUP
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Daemon is exiting.
				return
			case <-ch:
				// we got a SIGHUP. Rotate log file.
				if logger, ok := logWriter.(*lumberjack.Logger); ok {
					if err := logger.Rotate(); err != nil {
						slog.With(slog.Any("error", err)).ErrorContext(ctx, "failed to rotate log file")
					}
				}
			}
		}
	}()

	writers := []io.Writer{
		logWriter,
	}

	if hub != nil {
		writers = append(writers, hub)
	}

	wr := io.MultiWriter(writers...)

	var logLevel slog.Level
	switch logCfg.LogLevel {
	case config.LevelDebug:
		logLevel = slog.LevelDebug
	case config.LevelInfo:
		logLevel = slog.LevelInfo
	case config.LevelWarn:
		logLevel = slog.LevelWarn
	case config.LevelError:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// logger options
	opts := slog.HandlerOptions{
		AddSource: logCfg.LogSource,
		Level:     logLevel,
	}

	var han slog.Handler
	switch logCfg.LogFormat {
	case config.FormatJSON:
		han = slog.NewJSONHandler(wr, &opts)
	default:
		han = slog.NewTextHandler(wr, &opts)
	}

	wrapped := garmUtil.ContextHandler{
		Handler: han,
	}
	slog.SetDefault(slog.New(wrapped))
}

func maybeUpdateURLsFromConfig(cfg config.Config, store common.Store) error {
	info, err := store.ControllerInfo()
	if err != nil {
		return errors.Wrap(err, "fetching controller info")
	}

	var updateParams params.UpdateControllerParams

	if info.MetadataURL == "" && cfg.Default.MetadataURL != "" {
		updateParams.MetadataURL = &cfg.Default.MetadataURL
	}

	if info.CallbackURL == "" && cfg.Default.CallbackURL != "" {
		updateParams.CallbackURL = &cfg.Default.CallbackURL
	}

	if info.WebhookURL == "" && cfg.Default.WebhookURL != "" {
		updateParams.WebhookURL = &cfg.Default.WebhookURL
	}

	if updateParams.MetadataURL == nil && updateParams.CallbackURL == nil && updateParams.WebhookURL == nil {
		// nothing to update
		return nil
	}

	_, err = store.UpdateController(updateParams)
	if err != nil {
		return errors.Wrap(err, "updating controller info")
	}
	return nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()
	watcher.InitWatcher(ctx)

	ctx = auth.GetAdminContext(ctx)

	cfg, err := config.NewConfig(*conf)
	if err != nil {
		log.Fatalf("Fetching config: %+v", err) //nolint:gocritic
	}

	logCfg := cfg.GetLoggingConfig()
	var hub *websocket.Hub
	if logCfg.EnableLogStreamer != nil && *logCfg.EnableLogStreamer {
		hub = websocket.NewHub(ctx)
		if err := hub.Start(); err != nil {
			log.Fatal(err)
		}
		defer hub.Stop() //nolint
	}
	setupLogging(ctx, logCfg, hub)

	// Migrate credentials to the new format. This field will be read
	// by the DB migration logic.
	cfg.Database.MigrateCredentials = cfg.Github
	db, err := database.NewDatabase(ctx, cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	if err := maybeInitController(db); err != nil {
		log.Fatal(err)
	}

	if err := maybeUpdateURLsFromConfig(*cfg, db); err != nil {
		log.Fatal(err)
	}

	runner, err := runner.NewRunner(ctx, *cfg, db)
	if err != nil {
		log.Fatalf("failed to create controller: %+v", err)
	}

	// If there are many repos/pools, this may take a long time.
	if err := runner.Start(); err != nil {
		log.Fatal(err)
	}

	authenticator := auth.NewAuthenticator(cfg.JWTAuth, db)
	controller, err := controllers.NewAPIController(runner, authenticator, hub)
	if err != nil {
		log.Fatalf("failed to create controller: %+v", err)
	}

	instanceMiddleware, err := auth.NewInstanceMiddleware(db, cfg.JWTAuth)
	if err != nil {
		log.Fatal(err)
	}

	jwtMiddleware, err := auth.NewjwtMiddleware(db, cfg.JWTAuth)
	if err != nil {
		log.Fatal(err)
	}

	initMiddleware, err := auth.NewInitRequiredMiddleware(db)
	if err != nil {
		log.Fatal(err)
	}

	urlsRequiredMiddleware, err := auth.NewUrlsRequiredMiddleware(db)
	if err != nil {
		log.Fatal(err)
	}

	metricsMiddleware, err := auth.NewMetricsMiddleware(cfg.JWTAuth)
	if err != nil {
		log.Fatal(err)
	}

	router := routers.NewAPIRouter(controller, jwtMiddleware, initMiddleware, urlsRequiredMiddleware, instanceMiddleware, cfg.Default.EnableWebhookManagement)

	// start the metrics collector
	if cfg.Metrics.Enable {
		slog.InfoContext(ctx, "setting up metric routes")
		router = routers.WithMetricsRouter(router, cfg.Metrics.DisableAuth, metricsMiddleware)

		slog.InfoContext(ctx, "register metrics")
		if err := metrics.RegisterMetrics(); err != nil {
			log.Fatal(err)
		}

		slog.InfoContext(ctx, "start metrics collection")
		runnerMetrics.CollectObjectMetric(ctx, runner, cfg.Metrics.Duration())
	}

	if cfg.Default.DebugServer {
		slog.InfoContext(ctx, "setting up debug routes")
		router = routers.WithDebugServer(router)
	}

	corsMw := mux.CORSMethodMiddleware(router)
	router.Use(corsMw)

	allowedOrigins := handlers.AllowedOrigins(cfg.APIServer.CORSOrigins)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	// nolint:golangci-lint,gosec
	// G112: Potential Slowloris Attack because ReadHeaderTimeout is not configured in the http.Server
	srv := &http.Server{
		Addr: cfg.APIServer.BindAddress(),
		// Pass our instance of gorilla/mux in.
		Handler: handlers.CORS(methodsOk, headersOk, allowedOrigins)(router),
	}

	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf("creating listener: %q", err)
	}

	go func() {
		if cfg.APIServer.UseTLS {
			if err := srv.ServeTLS(listener, cfg.APIServer.TLSConfig.CRT, cfg.APIServer.TLSConfig.Key); err != nil {
				slog.With(slog.Any("error", err)).ErrorContext(ctx, "Listening")
			}
		} else {
			if err := srv.Serve(listener); err != http.ErrServerClosed {
				slog.With(slog.Any("error", err)).ErrorContext(ctx, "Listening")
			}
		}
	}()

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(ctx, "graceful api server shutdown failed")
	}

	slog.With(slog.Any("error", err)).ErrorContext(ctx, "waiting for runner to stop")
	if err := runner.Wait(); err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(ctx, "failed to shutdown workers")
		os.Exit(1)
	}
}
