# UpdateScaleSetParams


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **boolean** |  | [optional] [default to undefined]
**extended_state** | **string** |  | [optional] [default to undefined]
**extra_specs** | **object** |  | [optional] [default to undefined]
**flavor** | **string** |  | [optional] [default to undefined]
**image** | **string** |  | [optional] [default to undefined]
**max_runners** | **number** |  | [optional] [default to undefined]
**min_idle_runners** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**os_arch** | **string** |  | [optional] [default to undefined]
**os_type** | **string** |  | [optional] [default to undefined]
**runner_bootstrap_timeout** | **number** |  | [optional] [default to undefined]
**runner_group** | **string** | GithubRunnerGroup is the github runner group in which the runners of this pool will be added to. The runner group must be created by someone with access to the enterprise. | [optional] [default to undefined]
**runner_prefix** | **string** |  | [optional] [default to undefined]
**state** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { UpdateScaleSetParams } from './api';

const instance: UpdateScaleSetParams = {
    enabled,
    extended_state,
    extra_specs,
    flavor,
    image,
    max_runners,
    min_idle_runners,
    name,
    os_arch,
    os_type,
    runner_bootstrap_timeout,
    runner_group,
    runner_prefix,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
