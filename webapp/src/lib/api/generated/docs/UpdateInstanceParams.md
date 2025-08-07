# UpdateInstanceParams


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**addresses** | [**Array&lt;Address&gt;**](Address.md) | Addresses is a list of IP addresses the provider reports for this instance. | [optional] [default to undefined]
**os_name** | **string** | OSName is the name of the OS. Eg: ubuntu, centos, etc. | [optional] [default to undefined]
**os_version** | **string** | OSVersion is the version of the operating system. | [optional] [default to undefined]
**provider_fault** | **Array&lt;number&gt;** |  | [optional] [default to undefined]
**provider_id** | **string** |  | [optional] [default to undefined]
**runner_status** | **string** |  | [optional] [default to undefined]
**status** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { UpdateInstanceParams } from './api';

const instance: UpdateInstanceParams = {
    addresses,
    os_name,
    os_version,
    provider_fault,
    provider_id,
    runner_status,
    status,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
