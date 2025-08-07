# Instance


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**addresses** | [**Array&lt;Address&gt;**](Address.md) | Addresses is a list of IP addresses the provider reports for this instance. | [optional] [default to undefined]
**agent_id** | **number** | AgentID is the github runner agent ID. | [optional] [default to undefined]
**created_at** | **string** | CreatedAt is the timestamp of the creation of this runner. | [optional] [default to undefined]
**github_runner_group** | **string** | GithubRunnerGroup is the github runner group to which the runner belongs. The runner group must be created by someone with access to the enterprise. | [optional] [default to undefined]
**id** | **string** | ID is the database ID of this instance. | [optional] [default to undefined]
**job** | [**Job**](Job.md) |  | [optional] [default to undefined]
**name** | **string** | Name is the name associated with an instance. Depending on the provider, this may or may not be useful in the context of the provider, but we can use it internally to identify the instance. | [optional] [default to undefined]
**os_arch** | **string** |  | [optional] [default to undefined]
**os_name** | **string** | OSName is the name of the OS. Eg: ubuntu, centos, etc. | [optional] [default to undefined]
**os_type** | **string** |  | [optional] [default to undefined]
**os_version** | **string** | OSVersion is the version of the operating system. | [optional] [default to undefined]
**pool_id** | **string** | PoolID is the ID of the garm pool to which a runner belongs. | [optional] [default to undefined]
**provider_fault** | **Array&lt;number&gt;** | ProviderFault holds any error messages captured from the IaaS provider that is responsible for managing the lifecycle of the runner. | [optional] [default to undefined]
**provider_id** | **string** | PeoviderID is the unique ID the provider associated with the compute instance. We use this to identify the instance in the provider. | [optional] [default to undefined]
**provider_name** | **string** | ProviderName is the name of the IaaS where the instance was created. | [optional] [default to undefined]
**runner_status** | **string** |  | [optional] [default to undefined]
**scale_set_id** | **number** | ScaleSetID is the ID of the scale set to which a runner belongs. | [optional] [default to undefined]
**status** | **string** |  | [optional] [default to undefined]
**status_messages** | [**Array&lt;StatusMessage&gt;**](StatusMessage.md) | StatusMessages is a list of status messages sent back by the runner as it sets itself up. | [optional] [default to undefined]
**updated_at** | **string** | UpdatedAt is the timestamp of the last update to this runner. | [optional] [default to undefined]

## Example

```typescript
import { Instance } from './api';

const instance: Instance = {
    addresses,
    agent_id,
    created_at,
    github_runner_group,
    id,
    job,
    name,
    os_arch,
    os_name,
    os_type,
    os_version,
    pool_id,
    provider_fault,
    provider_id,
    provider_name,
    runner_status,
    scale_set_id,
    status,
    status_messages,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
