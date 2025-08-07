# Pool


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** |  | [optional] [default to undefined]
**enabled** | **boolean** |  | [optional] [default to undefined]
**endpoint** | [**ForgeEndpoint**](ForgeEndpoint.md) |  | [optional] [default to undefined]
**enterprise_id** | **string** |  | [optional] [default to undefined]
**enterprise_name** | **string** |  | [optional] [default to undefined]
**extra_specs** | **object** | ExtraSpecs is an opaque raw json that gets sent to the provider as part of the bootstrap params for instances. It can contain any kind of data needed by providers. The contents of this field means nothing to garm itself. We don\&#39;t act on the information in this field at all. We only validate that it\&#39;s a proper json. | [optional] [default to undefined]
**flavor** | **string** |  | [optional] [default to undefined]
**github_runner_group** | **string** | GithubRunnerGroup is the github runner group in which the runners will be added. The runner group must be created by someone with access to the enterprise. | [optional] [default to undefined]
**id** | **string** |  | [optional] [default to undefined]
**image** | **string** |  | [optional] [default to undefined]
**instances** | [**Array&lt;Instance&gt;**](Instance.md) |  | [optional] [default to undefined]
**max_runners** | **number** |  | [optional] [default to undefined]
**min_idle_runners** | **number** |  | [optional] [default to undefined]
**org_id** | **string** |  | [optional] [default to undefined]
**org_name** | **string** |  | [optional] [default to undefined]
**os_arch** | **string** |  | [optional] [default to undefined]
**os_type** | **string** |  | [optional] [default to undefined]
**priority** | **number** | Priority is the priority of the pool. The higher the number, the higher the priority. When fetching matching pools for a set of tags, the result will be sorted in descending order of priority. | [optional] [default to undefined]
**provider_name** | **string** |  | [optional] [default to undefined]
**repo_id** | **string** |  | [optional] [default to undefined]
**repo_name** | **string** |  | [optional] [default to undefined]
**runner_bootstrap_timeout** | **number** |  | [optional] [default to undefined]
**runner_prefix** | **string** |  | [optional] [default to undefined]
**tags** | [**Array&lt;Tag&gt;**](Tag.md) |  | [optional] [default to undefined]
**updated_at** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { Pool } from './api';

const instance: Pool = {
    created_at,
    enabled,
    endpoint,
    enterprise_id,
    enterprise_name,
    extra_specs,
    flavor,
    github_runner_group,
    id,
    image,
    instances,
    max_runners,
    min_idle_runners,
    org_id,
    org_name,
    os_arch,
    os_type,
    priority,
    provider_name,
    repo_id,
    repo_name,
    runner_bootstrap_timeout,
    runner_prefix,
    tags,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
