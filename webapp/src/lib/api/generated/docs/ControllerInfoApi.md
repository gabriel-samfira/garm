# ControllerInfoApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**controllerInfo**](#controllerinfo) | **GET** /controller-info | Get controller info.|

# **controllerInfo**
> ControllerInfo controllerInfo()


### Example

```typescript
import {
    ControllerInfoApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ControllerInfoApi(configuration);

const { status, data } = await apiInstance.controllerInfo();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ControllerInfo**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ControllerInfo |  -  |
|**409** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

