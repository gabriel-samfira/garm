# ProvidersApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**listProviders**](#listproviders) | **GET** /providers | List all providers.|

# **listProviders**
> Array<Provider> listProviders()


### Example

```typescript
import {
    ProvidersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ProvidersApi(configuration);

const { status, data } = await apiInstance.listProviders();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<Provider>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Providers |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

