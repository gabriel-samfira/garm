# MetricsTokenApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getMetricsToken**](#getmetricstoken) | **GET** /metrics-token | Returns a JWT token that can be used to access the metrics endpoint.|

# **getMetricsToken**
> JWTResponse getMetricsToken()


### Example

```typescript
import {
    MetricsTokenApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MetricsTokenApi(configuration);

const { status, data } = await apiInstance.getMetricsToken();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**JWTResponse**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | JWTResponse |  -  |
|**401** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

