# FirstRunApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**firstRun**](#firstrun) | **POST** /first-run | Initialize the first run of the controller.|

# **firstRun**
> User firstRun(body)


### Example

```typescript
import {
    FirstRunApi,
    Configuration,
    NewUserParams
} from './api';

const configuration = new Configuration();
const apiInstance = new FirstRunApi(configuration);

let body: NewUserParams; //Create a new user.

const { status, data } = await apiInstance.firstRun(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **NewUserParams**| Create a new user. | |


### Return type

**User**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | User |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

