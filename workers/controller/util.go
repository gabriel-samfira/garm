package controller

import (
	"fmt"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/params"
)

type entityGetter interface {
	GetEntity() (params.GithubEntity, error)
}

func githubEntityFromPayload(payload dbCommon.ChangePayload) (params.GithubEntity, error) {
	var entity entityGetter
	var ok bool
	switch payload.EntityType {
	case dbCommon.RepositoryEntityType:
		entity, ok = payload.Payload.(params.Repository)
	case dbCommon.OrganizationEntityType:
		entity, ok = payload.Payload.(params.Organization)
	case dbCommon.EnterpriseEntityType:
		entity, ok = payload.Payload.(params.Enterprise)
	default:
		return params.GithubEntity{}, fmt.Errorf("unsupported entity type %s", payload.EntityType)
	}
	if !ok {
		return params.GithubEntity{}, fmt.Errorf("payload is not a github entity type")
	}

	ret, err := entity.GetEntity()
	if err != nil {
		return params.GithubEntity{}, fmt.Errorf("failed to get entity: %w", err)
	}

	return ret, nil
}
