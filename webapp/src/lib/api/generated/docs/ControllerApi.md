# ControllerApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**updateController**](#updatecontroller) | **PUT** /controller | Update controller.|

# **updateController**
> ControllerInfo updateController(body)


### Example

```typescript
import {
    ControllerApi,
    Configuration,
    UpdateControllerParams
} from './api';

const configuration = new Configuration();
const apiInstance = new ControllerApi(configuration);

let body: UpdateControllerParams; //Parameters used when updating the controller.

const { status, data } = await apiInstance.updateController(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateControllerParams**| Parameters used when updating the controller. | |


### Return type

**ControllerInfo**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ControllerInfo |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

