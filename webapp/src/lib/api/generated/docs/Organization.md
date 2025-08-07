# Organization


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** |  | [optional] [default to undefined]
**credentials** | [**ForgeCredentials**](ForgeCredentials.md) |  | [optional] [default to undefined]
**credentials_id** | **number** |  | [optional] [default to undefined]
**credentials_name** | **string** | CredentialName is the name of the credentials associated with the enterprise. This field is now deprecated. Use CredentialsID instead. This field will be removed in v0.2.0. | [optional] [default to undefined]
**endpoint** | [**ForgeEndpoint**](ForgeEndpoint.md) |  | [optional] [default to undefined]
**events** | [**Array&lt;EntityEvent&gt;**](EntityEvent.md) |  | [optional] [default to undefined]
**id** | **string** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**pool** | [**Array&lt;Pool&gt;**](Pool.md) |  | [optional] [default to undefined]
**pool_balancing_type** | **string** |  | [optional] [default to undefined]
**pool_manager_status** | [**PoolManagerStatus**](PoolManagerStatus.md) |  | [optional] [default to undefined]
**updated_at** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { Organization } from './api';

const instance: Organization = {
    created_at,
    credentials,
    credentials_id,
    credentials_name,
    endpoint,
    events,
    id,
    name,
    pool,
    pool_balancing_type,
    pool_manager_status,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
