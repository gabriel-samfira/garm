# HooksApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getOrgWebhookInfo**](#getorgwebhookinfo) | **GET** /organizations/{orgID}/webhook | Get information about the GARM installed webhook on an organization.|
|[**getRepoWebhookInfo**](#getrepowebhookinfo) | **GET** /repositories/{repoID}/webhook | Get information about the GARM installed webhook on a repository.|
|[**installOrgWebhook**](#installorgwebhook) | **POST** /organizations/{orgID}/webhook | |
|[**installRepoWebhook**](#installrepowebhook) | **POST** /repositories/{repoID}/webhook | |
|[**uninstallOrgWebhook**](#uninstallorgwebhook) | **DELETE** /organizations/{orgID}/webhook | Uninstall organization webhook.|
|[**uninstallRepoWebhook**](#uninstallrepowebhook) | **DELETE** /repositories/{repoID}/webhook | Uninstall organization webhook.|

# **getOrgWebhookInfo**
> HookInfo getOrgWebhookInfo()


### Example

```typescript
import {
    HooksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new HooksApi(configuration);

let orgID: string; //Organization ID. (default to undefined)

const { status, data } = await apiInstance.getOrgWebhookInfo(
    orgID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


### Return type

**HookInfo**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | HookInfo |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRepoWebhookInfo**
> HookInfo getRepoWebhookInfo()


### Example

```typescript
import {
    HooksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new HooksApi(configuration);

let repoID: string; //Repository ID. (default to undefined)

const { status, data } = await apiInstance.getRepoWebhookInfo(
    repoID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**HookInfo**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | HookInfo |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **installOrgWebhook**
> HookInfo installOrgWebhook(body)

Install the GARM webhook for an organization. The secret configured on the organization will be used to validate the requests.

### Example

```typescript
import {
    HooksApi,
    Configuration,
    InstallWebhookParams
} from './api';

const configuration = new Configuration();
const apiInstance = new HooksApi(configuration);

let orgID: string; //Organization ID. (default to undefined)
let body: InstallWebhookParams; //Parameters used when creating the organization webhook.

const { status, data } = await apiInstance.installOrgWebhook(
    orgID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **InstallWebhookParams**| Parameters used when creating the organization webhook. | |
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


### Return type

**HookInfo**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | HookInfo |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **installRepoWebhook**
> HookInfo installRepoWebhook(body)

Install the GARM webhook for an organization. The secret configured on the organization will be used to validate the requests.

### Example

```typescript
import {
    HooksApi,
    Configuration,
    InstallWebhookParams
} from './api';

const configuration = new Configuration();
const apiInstance = new HooksApi(configuration);

let repoID: string; //Repository ID. (default to undefined)
let body: InstallWebhookParams; //Parameters used when creating the repository webhook.

const { status, data } = await apiInstance.installRepoWebhook(
    repoID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **InstallWebhookParams**| Parameters used when creating the repository webhook. | |
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**HookInfo**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | HookInfo |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uninstallOrgWebhook**
> APIErrorResponse uninstallOrgWebhook()


### Example

```typescript
import {
    HooksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new HooksApi(configuration);

let orgID: string; //Organization ID. (default to undefined)

const { status, data } = await apiInstance.uninstallOrgWebhook(
    orgID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


### Return type

**APIErrorResponse**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uninstallRepoWebhook**
> APIErrorResponse uninstallRepoWebhook()


### Example

```typescript
import {
    HooksApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new HooksApi(configuration);

let repoID: string; //Repository ID. (default to undefined)

const { status, data } = await apiInstance.uninstallRepoWebhook(
    repoID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**APIErrorResponse**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

