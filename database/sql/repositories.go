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

package sql

import (
	"context"
	"fmt"
	runnerErrors "garm/errors"
	"garm/params"
	"garm/util"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (s *sqlDatabase) CreateRepository(ctx context.Context, owner, name, credentialsName, webhookSecret string) (params.Repository, error) {
	secret := []byte{}
	var err error
	if webhookSecret != "" {
		secret, err = util.Aes256EncodeString(webhookSecret, s.cfg.Passphrase)
		if err != nil {
			return params.Repository{}, fmt.Errorf("failed to encrypt string")
		}
	}
	newRepo := Repository{
		Name:            name,
		Owner:           owner,
		WebhookSecret:   secret,
		CredentialsName: credentialsName,
	}

	q := s.conn.Create(&newRepo)
	if q.Error != nil {
		return params.Repository{}, errors.Wrap(q.Error, "creating repository")
	}

	param := s.sqlToCommonRepository(newRepo)
	param.WebhookSecret = webhookSecret

	return param, nil
}

func (s *sqlDatabase) GetRepository(ctx context.Context, owner, name string) (params.Repository, error) {
	repo, err := s.getRepo(ctx, owner, name)
	if err != nil {
		return params.Repository{}, errors.Wrap(err, "fetching repo")
	}

	param := s.sqlToCommonRepository(repo)
	secret, err := util.Aes256DecodeString(repo.WebhookSecret, s.cfg.Passphrase)
	if err != nil {
		return params.Repository{}, errors.Wrap(err, "decrypting secret")
	}
	param.WebhookSecret = secret

	return param, nil
}

func (s *sqlDatabase) ListRepositories(ctx context.Context) ([]params.Repository, error) {
	var repos []Repository
	q := s.conn.Find(&repos)
	if q.Error != nil {
		return []params.Repository{}, errors.Wrap(q.Error, "fetching user from database")
	}

	ret := make([]params.Repository, len(repos))
	for idx, val := range repos {
		ret[idx] = s.sqlToCommonRepository(val)
		if len(val.WebhookSecret) > 0 {
			secret, err := util.Aes256DecodeString(val.WebhookSecret, s.cfg.Passphrase)
			if err != nil {
				return nil, errors.Wrap(err, "decrypting secret")
			}
			ret[idx].WebhookSecret = secret
		}
	}

	return ret, nil
}

func (s *sqlDatabase) DeleteRepository(ctx context.Context, repoID string) error {
	repo, err := s.getRepoByID(ctx, repoID)
	if err != nil {
		return errors.Wrap(err, "fetching repo")
	}

	q := s.conn.Unscoped().Delete(&repo)
	if q.Error != nil && !errors.Is(q.Error, gorm.ErrRecordNotFound) {
		return errors.Wrap(q.Error, "deleting repo")
	}

	return nil
}

func (s *sqlDatabase) UpdateRepository(ctx context.Context, repoID string, param params.UpdateRepositoryParams) (params.Repository, error) {
	repo, err := s.getRepoByID(ctx, repoID)
	if err != nil {
		return params.Repository{}, errors.Wrap(err, "fetching repo")
	}

	if param.CredentialsName != "" {
		repo.CredentialsName = param.CredentialsName
	}

	if param.WebhookSecret != "" {
		secret, err := util.Aes256EncodeString(param.WebhookSecret, s.cfg.Passphrase)
		if err != nil {
			return params.Repository{}, fmt.Errorf("failed to encrypt string")
		}
		repo.WebhookSecret = secret
	}

	q := s.conn.Save(&repo)
	if q.Error != nil {
		return params.Repository{}, errors.Wrap(q.Error, "saving repo")
	}

	newParams := s.sqlToCommonRepository(repo)
	secret, err := util.Aes256DecodeString(repo.WebhookSecret, s.cfg.Passphrase)
	if err != nil {
		return params.Repository{}, errors.Wrap(err, "decrypting secret")
	}
	newParams.WebhookSecret = secret
	return newParams, nil
}

func (s *sqlDatabase) GetRepositoryByID(ctx context.Context, repoID string) (params.Repository, error) {
	repo, err := s.getRepoByID(ctx, repoID, "Pools")
	if err != nil {
		return params.Repository{}, errors.Wrap(err, "fetching repo")
	}

	param := s.sqlToCommonRepository(repo)
	secret, err := util.Aes256DecodeString(repo.WebhookSecret, s.cfg.Passphrase)
	if err != nil {
		return params.Repository{}, errors.Wrap(err, "decrypting secret")
	}
	param.WebhookSecret = secret

	return param, nil
}

func (s *sqlDatabase) CreateRepositoryPool(ctx context.Context, repoId string, param params.CreatePoolParams) (params.Pool, error) {
	if len(param.Tags) == 0 {
		return params.Pool{}, runnerErrors.NewBadRequestError("no tags specified")
	}

	repo, err := s.getRepoByID(ctx, repoId)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching repo")
	}

	newPool := Pool{
		ProviderName:           param.ProviderName,
		MaxRunners:             param.MaxRunners,
		MinIdleRunners:         param.MinIdleRunners,
		Image:                  param.Image,
		Flavor:                 param.Flavor,
		OSType:                 param.OSType,
		OSArch:                 param.OSArch,
		RepoID:                 repo.ID,
		Enabled:                param.Enabled,
		RunnerBootstrapTimeout: param.RunnerBootstrapTimeout,
	}

	_, err = s.getRepoPoolByUniqueFields(ctx, repoId, newPool.ProviderName, newPool.Image, newPool.Flavor)
	if err != nil {
		if !errors.Is(err, runnerErrors.ErrNotFound) {
			return params.Pool{}, errors.Wrap(err, "creating pool")
		}
	} else {
		return params.Pool{}, runnerErrors.NewConflictError("pool with the same image and flavor already exists on this provider")
	}

	tags := []Tag{}
	for _, val := range param.Tags {
		t, err := s.getOrCreateTag(val)
		if err != nil {
			return params.Pool{}, errors.Wrap(err, "fetching tag")
		}
		tags = append(tags, t)
	}

	q := s.conn.Create(&newPool)
	if q.Error != nil {
		return params.Pool{}, errors.Wrap(q.Error, "adding pool")
	}

	for _, tt := range tags {
		if err := s.conn.Model(&newPool).Association("Tags").Append(&tt); err != nil {
			return params.Pool{}, errors.Wrap(err, "saving tag")
		}
	}

	pool, err := s.getPoolByID(ctx, newPool.ID.String(), "Tags", "Instances", "Enterprise", "Organization", "Repository")
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool")
	}

	return s.sqlToCommonPool(pool), nil
}

func (s *sqlDatabase) ListRepoPools(ctx context.Context, repoID string) ([]params.Pool, error) {
	pools, err := s.getRepoPools(ctx, repoID, "Tags")
	if err != nil {
		return nil, errors.Wrap(err, "fetching pools")
	}

	ret := make([]params.Pool, len(pools))
	for idx, pool := range pools {
		ret[idx] = s.sqlToCommonPool(pool)
	}

	return ret, nil
}

func (s *sqlDatabase) GetRepositoryPool(ctx context.Context, repoID, poolID string) (params.Pool, error) {
	pool, err := s.getRepoPool(ctx, repoID, poolID, "Tags", "Instances")
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool")
	}
	return s.sqlToCommonPool(pool), nil
}

func (s *sqlDatabase) DeleteRepositoryPool(ctx context.Context, repoID, poolID string) error {
	pool, err := s.getRepoPool(ctx, repoID, poolID)
	if err != nil {
		return errors.Wrap(err, "looking up repo pool")
	}
	q := s.conn.Unscoped().Delete(&pool)
	if q.Error != nil && !errors.Is(q.Error, gorm.ErrRecordNotFound) {
		return errors.Wrap(q.Error, "deleting pool")
	}
	return nil
}

func (s *sqlDatabase) FindRepositoryPoolByTags(ctx context.Context, repoID string, tags []string) (params.Pool, error) {
	pool, err := s.findPoolByTags(repoID, "repo_id", tags)
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool")
	}
	return pool, nil
}

func (s *sqlDatabase) ListRepoInstances(ctx context.Context, repoID string) ([]params.Instance, error) {
	pools, err := s.getRepoPools(ctx, repoID, "Instances")
	if err != nil {
		return nil, errors.Wrap(err, "fetching repo")
	}

	ret := []params.Instance{}
	for _, pool := range pools {
		for _, instance := range pool.Instances {
			ret = append(ret, s.sqlToParamsInstance(instance))
		}
	}
	return ret, nil
}

func (s *sqlDatabase) UpdateRepositoryPool(ctx context.Context, repoID, poolID string, param params.UpdatePoolParams) (params.Pool, error) {
	pool, err := s.getRepoPool(ctx, repoID, poolID, "Tags", "Instances", "Enterprise", "Organization", "Repository")
	if err != nil {
		return params.Pool{}, errors.Wrap(err, "fetching pool")
	}

	return s.updatePool(pool, param)
}

func (s *sqlDatabase) getRepo(ctx context.Context, owner, name string) (Repository, error) {
	var repo Repository

	q := s.conn.Where("name = ? COLLATE NOCASE and owner = ? COLLATE NOCASE", name, owner).
		First(&repo)

	q = q.First(&repo)

	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			return Repository{}, runnerErrors.ErrNotFound
		}
		return Repository{}, errors.Wrap(q.Error, "fetching repository from database")
	}
	return repo, nil
}

func (s *sqlDatabase) findPoolByTags(id, poolType string, tags []string) (params.Pool, error) {
	if len(tags) == 0 {
		return params.Pool{}, runnerErrors.NewBadRequestError("missing tags")
	}
	u, err := uuid.FromString(id)
	if err != nil {
		return params.Pool{}, errors.Wrap(runnerErrors.ErrBadRequest, "parsing id")
	}

	var pool Pool
	where := fmt.Sprintf("tags.name in ? and %s = ?", poolType)
	q := s.conn.Joins("JOIN pool_tags on pool_tags.pool_id=pools.id").
		Joins("JOIN tags on tags.id=pool_tags.tag_id").
		Group("pools.id").
		Preload("Tags").
		Having("count(1) = ?", len(tags)).
		Where(where, tags, u).First(&pool)

	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			return params.Pool{}, runnerErrors.ErrNotFound
		}
		return params.Pool{}, errors.Wrap(q.Error, "fetching pool")
	}

	return s.sqlToCommonPool(pool), nil
}

func (s *sqlDatabase) getRepoPool(ctx context.Context, repoID, poolID string, preload ...string) (Pool, error) {
	repo, err := s.getRepoByID(ctx, repoID)
	if err != nil {
		return Pool{}, errors.Wrap(err, "fetching repo")
	}

	u, err := uuid.FromString(poolID)
	if err != nil {
		return Pool{}, errors.Wrap(runnerErrors.ErrBadRequest, "parsing id")
	}

	q := s.conn
	if len(preload) > 0 {
		for _, item := range preload {
			q = q.Preload(item)
		}
	}

	var pool []Pool
	err = q.Model(&repo).Association("Pools").Find(&pool, "id = ?", u)
	if err != nil {
		return Pool{}, errors.Wrap(err, "fetching pool")
	}
	if len(pool) == 0 {
		return Pool{}, runnerErrors.ErrNotFound
	}

	return pool[0], nil
}

func (s *sqlDatabase) getRepoPoolByUniqueFields(ctx context.Context, repoID string, provider, image, flavor string) (Pool, error) {
	repo, err := s.getRepoByID(ctx, repoID)
	if err != nil {
		return Pool{}, errors.Wrap(err, "fetching repo")
	}

	q := s.conn
	var pool []Pool
	err = q.Model(&repo).Association("Pools").Find(&pool, "provider_name = ? and image = ? and flavor = ?", provider, image, flavor)
	if err != nil {
		return Pool{}, errors.Wrap(err, "fetching pool")
	}
	if len(pool) == 0 {
		return Pool{}, runnerErrors.ErrNotFound
	}

	return pool[0], nil
}

func (s *sqlDatabase) getRepoPools(ctx context.Context, repoID string, preload ...string) ([]Pool, error) {
	repo, err := s.getRepoByID(ctx, repoID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching repo")
	}

	var pools []Pool
	q := s.conn.Model(&repo)
	if len(preload) > 0 {
		for _, item := range preload {
			q = q.Preload(item)
		}
	}
	err = q.Association("Pools").Find(&pools)
	if err != nil {
		return nil, errors.Wrap(err, "fetching pool")
	}

	return pools, nil
}

func (s *sqlDatabase) getRepoByID(ctx context.Context, id string, preload ...string) (Repository, error) {
	u, err := uuid.FromString(id)
	if err != nil {
		return Repository{}, errors.Wrap(runnerErrors.ErrBadRequest, "parsing id")
	}
	var repo Repository

	q := s.conn
	if len(preload) > 0 {
		for _, field := range preload {
			q = q.Preload(field)
		}
	}
	q = q.Where("id = ?", u).First(&repo)

	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			return Repository{}, runnerErrors.ErrNotFound
		}
		return Repository{}, errors.Wrap(q.Error, "fetching repository from database")
	}
	return repo, nil
}
