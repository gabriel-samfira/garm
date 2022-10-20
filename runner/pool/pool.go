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

package pool

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"garm/auth"
	dbCommon "garm/database/common"
	runnerErrors "garm/errors"
	"garm/params"
	"garm/runner/common"
	providerCommon "garm/runner/providers/common"

	"github.com/google/go-github/v48/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	poolIDLabelprefix     = "runner-pool-id:"
	controllerLabelPrefix = "runner-controller-id:"
)

const (
	// maxCreateAttempts is the number of times we will attempt to create an instance
	// before we give up.
	// TODO: make this configurable(?)
	maxCreateAttempts = 5
)

type basePool struct {
	ctx          context.Context
	controllerID string

	store dbCommon.Store

	providers map[string]common.Provider
	tools     []*github.RunnerApplicationDownload
	quit      chan struct{}
	done      chan struct{}

	helper       poolHelper
	credsDetails params.GithubCredentials

	mux sync.Mutex
}

func controllerIDFromLabels(labels []*github.RunnerLabels) string {
	for _, lbl := range labels {
		if lbl.Name != nil && strings.HasPrefix(*lbl.Name, controllerLabelPrefix) {
			labelName := *lbl.Name
			return labelName[len(controllerLabelPrefix):]
		}
	}
	return ""
}

// cleanupOrphanedProviderRunners compares runners in github with local runners and removes
// any local runners that are not present in Github. Runners that are "idle" in our
// provider, but do not exist in github, will be removed. This can happen if the
// garm was offline while a job was executed by a github action. When this
// happens, github will remove the ephemeral worker and send a webhook our way.
// If we were offline and did not process the webhook, the instance will linger.
// We need to remove it from the provider and database.
func (r *basePool) cleanupOrphanedProviderRunners(runners []*github.Runner) error {
	dbInstances, err := r.helper.FetchDbInstances()
	if err != nil {
		return errors.Wrap(err, "fetching instances from db")
	}

	runnerNames := map[string]bool{}
	for _, run := range runners {
		runnerNames[*run.Name] = true
	}

	for _, instance := range dbInstances {
		switch providerCommon.InstanceStatus(instance.Status) {
		case providerCommon.InstancePendingCreate,
			providerCommon.InstancePendingDelete:
			// this instance is in the process of being created or is awaiting deletion.
			// Instances in pending_create did not get a chance to register themselves in,
			// github so we let them be for now.
			continue
		}

		if ok := runnerNames[instance.Name]; !ok {
			// Set pending_delete on DB field. Allow consolidate() to remove it.
			if err := r.setInstanceStatus(instance.Name, providerCommon.InstancePendingDelete, nil); err != nil {
				log.Printf("failed to update runner %s status", instance.Name)
				return errors.Wrap(err, "updating runner")
			}
		}
	}
	return nil
}

// reapTimedOutRunners will mark as pending_delete any runner that has a status
// of "running" in the provider, but that has not registered with Github, and has
// received no new updates in the configured timeout interval.
func (r *basePool) reapTimedOutRunners(runners []*github.Runner) error {
	dbInstances, err := r.helper.FetchDbInstances()
	if err != nil {
		return errors.Wrap(err, "fetching instances from db")
	}

	runnerNames := map[string]bool{}
	for _, run := range runners {
		runnerNames[*run.Name] = true
	}

	for _, instance := range dbInstances {
		if ok := runnerNames[instance.Name]; !ok {
			if instance.Status == providerCommon.InstanceRunning {
				pool, err := r.store.GetPoolByID(r.ctx, instance.PoolID)
				if err != nil {
					return errors.Wrap(err, "fetching instance pool info")
				}
				if time.Since(instance.UpdatedAt).Minutes() < float64(pool.RunnerTimeout()) {
					continue
				}
				log.Printf("reaping instance %s due to timeout", instance.Name)
				if err := r.setInstanceStatus(instance.Name, providerCommon.InstancePendingDelete, nil); err != nil {
					log.Printf("failed to update runner %s status", instance.Name)
					return errors.Wrap(err, "updating runner")
				}
			}
		}
	}
	return nil
}

// cleanupOrphanedGithubRunners will forcefully remove any github runners that appear
// as offline and for which we no longer have a local instance.
// This may happen if someone manually deletes the instance in the provider. We need to
// first remove the instance from github, and then from our database.
func (r *basePool) cleanupOrphanedGithubRunners(runners []*github.Runner) error {
	for _, runner := range runners {
		runnerControllerID := controllerIDFromLabels(runner.Labels)
		if runnerControllerID != r.controllerID {
			// Not a runner we manage. Do not remove foreign runner.
			continue
		}

		status := runner.GetStatus()
		if status != "offline" {
			// Runner is online. Ignore it.
			continue
		}

		dbInstance, err := r.store.GetInstanceByName(r.ctx, *runner.Name)
		if err != nil {
			if !errors.Is(err, runnerErrors.ErrNotFound) {
				return errors.Wrap(err, "fetching instance from DB")
			}
			// We no longer have a DB entry for this instance, and the runner appears offline in github.
			// Previous forceful removal may have failed?
			log.Printf("Runner %s has no database entry in garm, removing from github", *runner.Name)
			resp, err := r.helper.RemoveGithubRunner(*runner.ID)
			if err != nil {
				// Removed in the meantime?
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					continue
				}
				return errors.Wrap(err, "removing runner")
			}
			continue
		}

		if providerCommon.InstanceStatus(dbInstance.Status) == providerCommon.InstancePendingDelete {
			// already marked for deleting, which means the github workflow finished.
			// Let consolidate take care of it.
			continue
		}

		pool, err := r.helper.GetPoolByID(dbInstance.PoolID)
		if err != nil {
			return errors.Wrap(err, "fetching pool")
		}

		// check if the provider still has the instance.
		provider, ok := r.providers[pool.ProviderName]
		if !ok {
			return fmt.Errorf("unknown provider %s for pool %s", pool.ProviderName, pool.ID)
		}

		// Check if the instance is still on the provider.
		_, err = provider.GetInstance(r.ctx, dbInstance.Name)
		if err != nil {
			if !errors.Is(err, runnerErrors.ErrNotFound) {
				return errors.Wrap(err, "fetching instance from provider")
			}
			// The runner instance is no longer on the provider, and it appears offline in github.
			// It should be safe to force remove it.
			log.Printf("Runner instance for %s is no longer on the provider, removing from github", dbInstance.Name)
			resp, err := r.helper.RemoveGithubRunner(*runner.ID)
			if err != nil {
				// Removed in the meantime?
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					log.Printf("runner dissapeared from github")
				} else {
					return errors.Wrap(err, "removing runner from github")
				}
			}
			// Remove the database entry for the runner.
			log.Printf("Removing %s from database", dbInstance.Name)
			if err := r.store.DeleteInstance(r.ctx, dbInstance.PoolID, dbInstance.Name); err != nil {
				return errors.Wrap(err, "removing runner from database")
			}
			continue
		}

		if providerCommon.InstanceStatus(dbInstance.Status) == providerCommon.InstanceRunning {
			// instance is running, but github reports runner as offline. Log the event.
			// This scenario requires manual intervention.
			// Perhaps it just came online and github did not yet change it's status?
			log.Printf("instance %s is online but github reports runner as offline", dbInstance.Name)
			continue
		}
		//start the instance
		if err := provider.Start(r.ctx, dbInstance.ProviderID); err != nil {
			return errors.Wrapf(err, "starting instance %s", dbInstance.ProviderID)
		}
	}
	return nil
}

func (r *basePool) fetchInstance(runnerName string) (params.Instance, error) {
	runner, err := r.store.GetInstanceByName(r.ctx, runnerName)
	if err != nil {
		return params.Instance{}, errors.Wrap(err, "fetching instance")
	}

	return runner, nil
}

func (r *basePool) setInstanceRunnerStatus(runnerName string, status providerCommon.RunnerStatus) error {
	updateParams := params.UpdateInstanceParams{
		RunnerStatus: status,
	}

	if err := r.updateInstance(runnerName, updateParams); err != nil {
		return errors.Wrap(err, "updating runner state")
	}
	return nil
}

func (r *basePool) updateInstance(runnerName string, update params.UpdateInstanceParams) error {
	runner, err := r.fetchInstance(runnerName)
	if err != nil {
		return errors.Wrap(err, "fetching instance")
	}

	if _, err := r.store.UpdateInstance(r.ctx, runner.ID, update); err != nil {
		return errors.Wrap(err, "updating runner state")
	}
	return nil
}

func (r *basePool) setInstanceStatus(runnerName string, status providerCommon.InstanceStatus, providerFault []byte) error {
	updateParams := params.UpdateInstanceParams{
		Status:        status,
		ProviderFault: providerFault,
	}

	if err := r.updateInstance(runnerName, updateParams); err != nil {
		return errors.Wrap(err, "updating runner state")
	}
	return nil
}

func (r *basePool) acquireNewInstance(job params.WorkflowJob) error {
	requestedLabels := job.WorkflowJob.Labels
	if len(requestedLabels) == 0 {
		// no labels were requested.
		return nil
	}

	pool, err := r.helper.FindPoolByTags(requestedLabels)
	if err != nil {
		return errors.Wrap(err, "fetching suitable pool")
	}
	log.Printf("adding new runner with requested tags %s in pool %s", strings.Join(job.WorkflowJob.Labels, ", "), pool.ID)

	if !pool.Enabled {
		log.Printf("selected pool (%s) is disabled", pool.ID)
		return nil
	}

	poolInstances, err := r.store.PoolInstanceCount(r.ctx, pool.ID)
	if err != nil {
		return errors.Wrap(err, "fetching instances")
	}

	if poolInstances >= int64(pool.MaxRunners) {
		log.Printf("max_runners (%d) reached for pool %s, skipping...", pool.MaxRunners, pool.ID)
		return nil
	}

	if err := r.AddRunner(r.ctx, pool.ID); err != nil {
		log.Printf("failed to add runner to pool %s", pool.ID)
		return errors.Wrap(err, "adding runner")
	}
	return nil
}

func (r *basePool) AddRunner(ctx context.Context, poolID string) error {
	pool, err := r.helper.GetPoolByID(poolID)
	if err != nil {
		return errors.Wrap(err, "fetching pool")
	}

	name := fmt.Sprintf("garm-%s", uuid.New())

	createParams := params.CreateInstanceParams{
		Name:          name,
		Status:        providerCommon.InstancePendingCreate,
		RunnerStatus:  providerCommon.RunnerPending,
		OSArch:        pool.OSArch,
		OSType:        pool.OSType,
		CallbackURL:   r.helper.GetCallbackURL(),
		CreateAttempt: 1,
	}

	_, err = r.store.CreateInstance(r.ctx, poolID, createParams)
	if err != nil {
		return errors.Wrap(err, "creating instance")
	}

	return nil
}

func (r *basePool) loop() {
	consolidateTimer := time.NewTicker(common.PoolConsilitationInterval)
	reapTimer := time.NewTicker(common.PoolReapTimeoutInterval)
	toolUpdateTimer := time.NewTicker(common.PoolToolUpdateInterval)
	defer func() {
		log.Printf("%s loop exited", r.helper.String())
		consolidateTimer.Stop()
		reapTimer.Stop()
		toolUpdateTimer.Stop()
		close(r.done)
	}()
	log.Printf("starting loop for %s", r.helper.String())
	// TODO: Consolidate runners on loop start. Provider runners must match runners
	// in github and DB. When a Workflow job is received, we will first create/update
	// an entity in the database, before sending the request to the provider to create/delete
	// an instance. If a "queued" job is received, we create an entity in the db with
	// a state of "pending_create". Once that instance is up and calls home, it is marked
	// as "active". If a "completed" job is received from github, we mark the instance
	// as "pending_delete". Once the provider deletes the instance, we mark it as "deleted"
	// in the database.
	// We also ensure we have runners created based on pool characteristics. This is where
	// we spin up "MinWorkers" for each runner type.

	for {
		select {
		case <-reapTimer.C:
			runners, err := r.helper.GetGithubRunners()
			if err != nil {
				log.Printf("error fetching github runners: %s", err)
				continue
			}
			if err := r.reapTimedOutRunners(runners); err != nil {
				log.Printf("failed to reap timed out runners: %q", err)
			}

			if err := r.cleanupOrphanedGithubRunners(runners); err != nil {
				log.Printf("failed to clean orphaned github runners: %q", err)
			}
		case <-consolidateTimer.C:
			// consolidate.
			r.consolidate()
		case <-toolUpdateTimer.C:
			// Update tools cache.
			tools, err := r.helper.FetchTools()
			if err != nil {
				log.Printf("failed to update tools for repo %s: %s", r.helper.String(), err)
			}
			r.mux.Lock()
			r.tools = tools
			r.mux.Unlock()
		case <-r.ctx.Done():
			// daemon is shutting down.
			return
		case <-r.quit:
			// this worker was stopped.
			return
		}
	}
}

func (r *basePool) addInstanceToProvider(instance params.Instance) error {
	pool, err := r.helper.GetPoolByID(instance.PoolID)
	if err != nil {
		return errors.Wrap(err, "fetching pool")
	}

	provider, ok := r.providers[pool.ProviderName]
	if !ok {
		return runnerErrors.NewNotFoundError("invalid provider ID")
	}

	labels := []string{}
	for _, tag := range pool.Tags {
		labels = append(labels, tag.Name)
	}
	labels = append(labels, r.controllerLabel())
	labels = append(labels, r.poolLabel(pool.ID))

	tk, err := r.helper.GetGithubRegistrationToken()
	if err != nil {
		return errors.Wrap(err, "fetching registration token")
	}

	jwtValidity := pool.RunnerTimeout()
	var poolType common.PoolType = common.RepositoryPool
	if pool.OrgID != "" {
		poolType = common.OrganizationPool
	}

	entity := r.helper.String()
	jwtToken, err := auth.NewInstanceJWTToken(instance, r.helper.JwtToken(), entity, poolType, jwtValidity)
	if err != nil {
		return errors.Wrap(err, "fetching instance jwt token")
	}

	bootstrapArgs := params.BootstrapInstance{
		Name:                    instance.Name,
		Tools:                   r.tools,
		RepoURL:                 r.helper.GithubURL(),
		GithubRunnerAccessToken: tk,
		CallbackURL:             instance.CallbackURL,
		InstanceToken:           jwtToken,
		OSArch:                  pool.OSArch,
		Flavor:                  pool.Flavor,
		Image:                   pool.Image,
		Labels:                  labels,
		PoolID:                  instance.PoolID,
		CACertBundle:            r.credsDetails.CABundle,
	}

	var instanceIDToDelete string

	defer func() {
		if instanceIDToDelete != "" {
			if err := provider.DeleteInstance(r.ctx, instanceIDToDelete); err != nil {
				if !errors.Is(err, runnerErrors.ErrNotFound) {
					log.Printf("failed to cleanup instance: %s", instanceIDToDelete)
				}
			}
		}
	}()

	providerInstance, err := provider.CreateInstance(r.ctx, bootstrapArgs)
	if err != nil {
		instanceIDToDelete = instance.Name
		return errors.Wrap(err, "creating instance")
	}

	if providerInstance.Status == providerCommon.InstanceError {
		instanceIDToDelete = instance.ProviderID
		if instanceIDToDelete == "" {
			instanceIDToDelete = instance.Name
		}
	}

	updateInstanceArgs := r.updateArgsFromProviderInstance(providerInstance)
	if _, err := r.store.UpdateInstance(r.ctx, instance.ID, updateInstanceArgs); err != nil {
		return errors.Wrap(err, "updating instance")
	}
	return nil
}

func (r *basePool) getRunnerNameFromJob(job params.WorkflowJob) (string, error) {
	if job.WorkflowJob.RunnerName != "" {
		return job.WorkflowJob.RunnerName, nil
	}

	// Runner name was not set in WorkflowJob by github. We can still attempt to
	// fetch the info we need, using the workflow run ID, from the API.
	log.Printf("runner name not found in workflow job, attempting to fetch from API")
	runnerName, err := r.helper.GetRunnerNameFromWorkflow(job)
	if err != nil {
		return "", errors.Wrap(err, "fetching runner name from API")
	}

	return runnerName, nil
}

func (r *basePool) HandleWorkflowJob(job params.WorkflowJob) error {
	if err := r.helper.ValidateOwner(job); err != nil {
		return errors.Wrap(err, "validating owner")
	}

	switch job.Action {
	case "queued":
		// Create instance in database and set it to pending create.
		// If we already have an idle runner around, that runner will pick up the job
		// and trigger an "in_progress" update from github (see bellow), which in turn will set the
		// runner state of the instance to "active". The ensureMinIdleRunners() function will
		// exclude that runner from available runners and attempt to ensure
		// the needed number of runners.
		if err := r.acquireNewInstance(job); err != nil {
			log.Printf("failed to add instance: %s", err)
		}
	case "completed":
		// ignore the error here. A completed job may not have a runner name set
		// if it was never assigned to a runner, and was canceled.
		runnerName, _ := r.getRunnerNameFromJob(job)
		// Set instance in database to pending delete.
		if runnerName == "" {
			// Unassigned jobs will have an empty runner_name.
			// There is nothing to to in this case.
			log.Printf("no runner was assigned. Skipping.")
			return nil
		}
		// update instance workload state.
		if err := r.setInstanceRunnerStatus(runnerName, providerCommon.RunnerTerminated); err != nil {
			log.Printf("failed to update runner %s status", runnerName)
			return errors.Wrap(err, "updating runner")
		}
		log.Printf("marking instance %s as pending_delete", runnerName)
		if err := r.setInstanceStatus(runnerName, providerCommon.InstancePendingDelete, nil); err != nil {
			log.Printf("failed to update runner %s status", runnerName)
			return errors.Wrap(err, "updating runner")
		}
	case "in_progress":
		// in_progress jobs must have a runner name/ID assigned. Sometimes github will send a hook without
		// a runner set. In such cases, we attemt to fetch it from the API.
		runnerName, err := r.getRunnerNameFromJob(job)
		if err != nil {
			return errors.Wrap(err, "determining runner name")
		}
		// update instance workload state.
		if err := r.setInstanceRunnerStatus(runnerName, providerCommon.RunnerActive); err != nil {
			log.Printf("failed to update runner %s status", job.WorkflowJob.RunnerName)
			return errors.Wrap(err, "updating runner")
		}
	}
	return nil
}

func (r *basePool) poolLabel(poolID string) string {
	return fmt.Sprintf("%s%s", poolIDLabelprefix, poolID)
}

func (r *basePool) controllerLabel() string {
	return fmt.Sprintf("%s%s", controllerLabelPrefix, r.controllerID)
}

func (r *basePool) updateArgsFromProviderInstance(providerInstance params.Instance) params.UpdateInstanceParams {
	return params.UpdateInstanceParams{
		ProviderID:    providerInstance.ProviderID,
		OSName:        providerInstance.OSName,
		OSVersion:     providerInstance.OSVersion,
		Addresses:     providerInstance.Addresses,
		Status:        providerInstance.Status,
		RunnerStatus:  providerInstance.RunnerStatus,
		ProviderFault: providerInstance.ProviderFault,
	}
}

func (r *basePool) ensureIdleRunnersForOnePool(pool params.Pool) {
	if !pool.Enabled {
		return
	}
	existingInstances, err := r.store.ListPoolInstances(r.ctx, pool.ID)
	if err != nil {
		log.Printf("failed to ensure minimum idle workers for pool %s: %s", pool.ID, err)
		return
	}

	if uint(len(existingInstances)) >= pool.MaxRunners {
		log.Printf("max workers (%d) reached for pool %s, skipping idle worker creation", pool.MaxRunners, pool.ID)
		return
	}

	idleOrPendingWorkers := []params.Instance{}
	for _, inst := range existingInstances {
		if providerCommon.RunnerStatus(inst.RunnerStatus) != providerCommon.RunnerActive {
			idleOrPendingWorkers = append(idleOrPendingWorkers, inst)
		}
	}

	var required int
	if len(idleOrPendingWorkers) < int(pool.MinIdleRunners) {
		// get the needed delta.
		required = int(pool.MinIdleRunners) - len(idleOrPendingWorkers)

		projectedInstanceCount := len(existingInstances) + required
		if uint(projectedInstanceCount) > pool.MaxRunners {
			// ensure we don't go above max workers
			delta := projectedInstanceCount - int(pool.MaxRunners)
			required = required - delta
		}
	}

	for i := 0; i < required; i++ {
		log.Printf("addind new idle worker to pool %s", pool.ID)
		if err := r.AddRunner(r.ctx, pool.ID); err != nil {
			log.Printf("failed to add new instance for pool %s: %s", pool.ID, err)
		}
	}
}

func (r *basePool) retryFailedInstancesForOnePool(pool params.Pool) {
	if !pool.Enabled {
		return
	}

	existingInstances, err := r.store.ListPoolInstances(r.ctx, pool.ID)
	if err != nil {
		log.Printf("failed to ensure minimum idle workers for pool %s: %s", pool.ID, err)
		return
	}

	for _, instance := range existingInstances {
		if instance.Status != providerCommon.InstanceError {
			continue
		}
		if instance.CreateAttempt >= maxCreateAttempts {
			continue
		}

		// NOTE(gabriel-samfira): this is done in parallel. If there are many failed instances
		// this has the potential to create many API requests to the target provider.
		// TODO(gabriel-samfira): implement request throttling.
		if instance.ProviderID == "" && instance.Name == "" {
			// This really should not happen, but no harm in being extra paranoid. The name is set
			// when creating a db entity for the runner, so we should at least have a name.
			return
		}
		// TODO(gabriel-samfira): Incrementing CreateAttempt should be done within a transaction.
		// It's fairly safe to do here (for now), as there should be no other code path that updates
		// an instance in this state.
		updateParams := params.UpdateInstanceParams{
			CreateAttempt: instance.CreateAttempt + 1,
			Status:        providerCommon.InstancePendingCreate,
		}
		log.Printf("queueing previously failed instance %s for retry", instance.Name)
		// Set instance to pending create and wait for retry.
		if err := r.updateInstance(instance.Name, updateParams); err != nil {
			log.Printf("failed to update runner %s status", instance.Name)
		}
	}
}

func (r *basePool) retryFailedInstances() {
	pools, err := r.helper.ListPools()
	if err != nil {
		log.Printf("error listing pools: %s", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(pools))
	for _, pool := range pools {
		go func(pool params.Pool) {
			defer wg.Done()
			r.retryFailedInstancesForOnePool(pool)
		}(pool)
	}
	wg.Wait()
}

func (r *basePool) ensureMinIdleRunners() {
	pools, err := r.helper.ListPools()
	if err != nil {
		log.Printf("error listing pools: %s", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(pools))
	for _, pool := range pools {
		go func(pool params.Pool) {
			defer wg.Done()
			r.ensureIdleRunnersForOnePool(pool)
		}(pool)
	}
	wg.Wait()
}

func (r *basePool) deleteInstanceFromProvider(instance params.Instance) error {
	pool, err := r.helper.GetPoolByID(instance.PoolID)
	if err != nil {
		return errors.Wrap(err, "fetching pool")
	}

	provider, ok := r.providers[pool.ProviderName]
	if !ok {
		return runnerErrors.NewNotFoundError("invalid provider ID")
	}

	identifier := instance.ProviderID
	if identifier == "" {
		// provider did not return a provider ID?
		// try with name
		identifier = instance.Name
	}

	if err := provider.DeleteInstance(r.ctx, identifier); err != nil {
		return errors.Wrap(err, "removing instance")
	}

	if err := r.store.DeleteInstance(r.ctx, pool.ID, instance.Name); err != nil {
		return errors.Wrap(err, "deleting instance from database")
	}
	return nil
}

func (r *basePool) deletePendingInstances() {
	instances, err := r.helper.FetchDbInstances()
	if err != nil {
		log.Printf("failed to fetch instances from store: %s", err)
		return
	}

	for _, instance := range instances {
		if instance.Status != providerCommon.InstancePendingDelete {
			// not in pending_delete status. Skip.
			continue
		}

		// Set the status to deleting before launching the goroutine that removes
		// the runner from the provider (which can take a long time).
		if err := r.setInstanceStatus(instance.Name, providerCommon.InstanceDeleting, nil); err != nil {
			log.Printf("failed to update runner %s status", instance.Name)
		}
		go func(instance params.Instance) (err error) {
			defer func(instance params.Instance) {
				if err != nil {
					// failed to remove from provider. Set the status back to pending_delete, which
					// will retry the operation.
					if err := r.setInstanceStatus(instance.Name, providerCommon.InstancePendingDelete, nil); err != nil {
						log.Printf("failed to update runner %s status", instance.Name)
					}
				}
			}(instance)

			err = r.deleteInstanceFromProvider(instance)
			if err != nil {
				log.Printf("failed to delete instance from provider: %+v", err)
			}
			return
		}(instance)
	}
}

func (r *basePool) addPendingInstances() {
	// TODO: filter instances by status.
	instances, err := r.helper.FetchDbInstances()
	if err != nil {
		log.Printf("failed to fetch instances from store: %s", err)
		return
	}

	for _, instance := range instances {
		if instance.Status != providerCommon.InstancePendingCreate {
			// not in pending_create status. Skip.
			continue
		}
		// Set the instance to "creating" before launching the goroutine. This will ensure that addPendingInstances()
		// won't attempt to create the runner a second time.
		if err := r.setInstanceStatus(instance.Name, providerCommon.InstanceCreating, nil); err != nil {
			log.Printf("failed to update runner %s status", instance.Name)
		}
		go func(instance params.Instance) {
			log.Printf("creating instance %s in pool %s", instance.Name, instance.PoolID)
			if err := r.addInstanceToProvider(instance); err != nil {
				log.Printf("failed to add instance to provider: %s", err)
				errAsBytes := []byte(err.Error())
				if err := r.setInstanceStatus(instance.Name, providerCommon.InstanceError, errAsBytes); err != nil {
					log.Printf("failed to update runner %s status", instance.Name)
				}
				log.Printf("failed to create instance in provider: %s", err)
			}
		}(instance)
	}
}

func (r *basePool) consolidate() {
	// TODO(gabriel-samfira): replace this with something more efficient.
	r.mux.Lock()
	defer r.mux.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		r.deletePendingInstances()
	}()
	go func() {
		defer wg.Done()
		r.addPendingInstances()
	}()
	wg.Wait()

	wg.Add(2)
	go func() {
		defer wg.Done()
		r.ensureMinIdleRunners()
	}()

	go func() {
		defer wg.Done()
		r.retryFailedInstances()
	}()
	wg.Wait()
}

func (r *basePool) Wait() error {
	select {
	case <-r.done:
	case <-time.After(20 * time.Second):
		return errors.Wrap(runnerErrors.ErrTimeout, "waiting for pool to stop")
	}
	return nil
}

func (r *basePool) Start() error {
	tools, err := r.helper.FetchTools()
	if err != nil {
		return errors.Wrap(err, "initializing tools")
	}
	r.mux.Lock()
	r.tools = tools
	r.mux.Unlock()

	runners, err := r.helper.GetGithubRunners()
	if err != nil {
		return errors.Wrap(err, "fetching github runners")
	}
	if err := r.cleanupOrphanedProviderRunners(runners); err != nil {
		return errors.Wrap(err, "cleaning orphaned instances")
	}

	if err := r.cleanupOrphanedGithubRunners(runners); err != nil {
		return errors.Wrap(err, "cleaning orphaned github runners")
	}
	go r.loop()
	return nil
}

func (r *basePool) Stop() error {
	close(r.quit)
	return nil
}

func (r *basePool) RefreshState(param params.UpdatePoolStateParams) error {
	return r.helper.UpdateState(param)
}

func (r *basePool) WebhookSecret() string {
	return r.helper.WebhookSecret()
}

func (r *basePool) ID() string {
	return r.helper.ID()
}

func (r *basePool) ForceDeleteRunner(runner params.Instance) error {
	if runner.AgentID != 0 {
		resp, err := r.helper.RemoveGithubRunner(runner.AgentID)
		if err != nil {
			if resp != nil {
				switch resp.StatusCode {
				case http.StatusUnprocessableEntity:
					return errors.Wrapf(runnerErrors.ErrBadRequest, "removing runner: %q", err)
				case http.StatusNotFound:
					// Runner may have been deleted by a finished job, or manually by the user.
					log.Printf("runner with agent id %d was not found in github", runner.AgentID)
				default:
					return errors.Wrap(err, "removing runner")
				}
			} else {
				// We got a nil response. Assume we are in error.
				return errors.Wrap(err, "removing runner")
			}
		}
	}
	log.Printf("setting instance status for: %v", runner.Name)

	if err := r.setInstanceStatus(runner.Name, providerCommon.InstancePendingDelete, nil); err != nil {
		log.Printf("failed to update runner %s status", runner.Name)
		return errors.Wrap(err, "updating runner")
	}
	return nil
}
