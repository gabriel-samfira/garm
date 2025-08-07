# LoginApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**login**](#login) | **POST** /auth/login | Logs in a user and returns a JWT token.|

# **login**
> JWTResponse login(body)


### Example

```typescript
import {
    LoginApi,
    Configuration,
    PasswordLoginParams
} from './api';

const configuration = new Configuration();
const apiInstance = new LoginApi(configuration);

let body: PasswordLoginParams; //Login information.

const { status, data } = await apiInstance.login(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **PasswordLoginParams**| Login information. | |


### Return type

**JWTResponse**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | JWTResponse |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

