# ForgeCredentials


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**api_base_url** | **string** |  | [optional] [default to undefined]
**auth_type** | **string** |  | [optional] [default to undefined]
**base_url** | **string** |  | [optional] [default to undefined]
**ca_bundle** | **Array&lt;number&gt;** |  | [optional] [default to undefined]
**created_at** | **string** |  | [optional] [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**endpoint** | [**ForgeEndpoint**](ForgeEndpoint.md) |  | [optional] [default to undefined]
**enterprises** | [**Array&lt;Enterprise&gt;**](Enterprise.md) |  | [optional] [default to undefined]
**forge_type** | **string** |  | [optional] [default to undefined]
**id** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**organizations** | [**Array&lt;Organization&gt;**](Organization.md) |  | [optional] [default to undefined]
**rate_limit** | [**GithubRateLimit**](GithubRateLimit.md) |  | [optional] [default to undefined]
**repositories** | [**Array&lt;Repository&gt;**](Repository.md) |  | [optional] [default to undefined]
**updated_at** | **string** |  | [optional] [default to undefined]
**upload_base_url** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { ForgeCredentials } from './api';

const instance: ForgeCredentials = {
    api_base_url,
    auth_type,
    base_url,
    ca_bundle,
    created_at,
    description,
    endpoint,
    enterprises,
    forge_type,
    id,
    name,
    organizations,
    rate_limit,
    repositories,
    updated_at,
    upload_base_url,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
