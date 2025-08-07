# CreateScaleSetParams


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**disable_update** | **boolean** |  | [optional] [default to undefined]
**enabled** | **boolean** |  | [optional] [default to undefined]
**extra_specs** | **object** |  | [optional] [default to undefined]
**flavor** | **string** |  | [optional] [default to undefined]
**github_runner_group** | **string** | GithubRunnerGroup is the github runner group in which the runners of this pool will be added to. The runner group must be created by someone with access to the enterprise. | [optional] [default to undefined]
**image** | **string** |  | [optional] [default to undefined]
**max_runners** | **number** |  | [optional] [default to undefined]
**min_idle_runners** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**os_arch** | **string** |  | [optional] [default to undefined]
**os_type** | **string** |  | [optional] [default to undefined]
**provider_name** | **string** |  | [optional] [default to undefined]
**runner_bootstrap_timeout** | **number** |  | [optional] [default to undefined]
**runner_prefix** | **string** |  | [optional] [default to undefined]
**scale_set_id** | **number** |  | [optional] [default to undefined]
**tags** | **Array&lt;string&gt;** |  | [optional] [default to undefined]

## Example

```typescript
import { CreateScaleSetParams } from './api';

const instance: CreateScaleSetParams = {
    disable_update,
    enabled,
    extra_specs,
    flavor,
    github_runner_group,
    image,
    max_runners,
    min_idle_runners,
    name,
    os_arch,
    os_type,
    provider_name,
    runner_bootstrap_timeout,
    runner_prefix,
    scale_set_id,
    tags,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
