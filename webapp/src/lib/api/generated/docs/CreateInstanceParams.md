# CreateInstanceParams


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**aditional_labels** | **Array&lt;string&gt;** |  | [optional] [default to undefined]
**callback_url** | **string** |  | [optional] [default to undefined]
**github_runner_group** | **string** | GithubRunnerGroup is the github runner group to which the runner belongs. The runner group must be created by someone with access to the enterprise. | [optional] [default to undefined]
**jit_configuration** | **{ [key: string]: string; }** |  | [optional] [default to undefined]
**metadata_url** | **string** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**os_arch** | **string** |  | [optional] [default to undefined]
**os_type** | **string** |  | [optional] [default to undefined]
**runner_status** | **string** |  | [optional] [default to undefined]
**status** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { CreateInstanceParams } from './api';

const instance: CreateInstanceParams = {
    aditional_labels,
    callback_url,
    github_runner_group,
    jit_configuration,
    metadata_url,
    name,
    os_arch,
    os_type,
    runner_status,
    status,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
