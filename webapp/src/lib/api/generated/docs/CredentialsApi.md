# CredentialsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createCredentials**](#createcredentials) | **POST** /github/credentials | Create a GitHub credential.|
|[**createGiteaCredentials**](#creategiteacredentials) | **POST** /gitea/credentials | Create a Gitea credential.|
|[**deleteCredentials**](#deletecredentials) | **DELETE** /github/credentials/{id} | Delete a GitHub credential.|
|[**deleteGiteaCredentials**](#deletegiteacredentials) | **DELETE** /gitea/credentials/{id} | Delete a Gitea credential.|
|[**getCredentials**](#getcredentials) | **GET** /github/credentials/{id} | Get a GitHub credential.|
|[**getGiteaCredentials**](#getgiteacredentials) | **GET** /gitea/credentials/{id} | Get a Gitea credential.|
|[**listCredentials**](#listcredentials) | **GET** /github/credentials | List all credentials.|
|[**listGiteaCredentials**](#listgiteacredentials) | **GET** /gitea/credentials | List all credentials.|
|[**updateCredentials**](#updatecredentials) | **PUT** /github/credentials/{id} | Update a GitHub credential.|
|[**updateGiteaCredentials**](#updategiteacredentials) | **PUT** /gitea/credentials/{id} | Update a Gitea credential.|

# **createCredentials**
> ForgeCredentials createCredentials(body)


### Example

```typescript
import {
    CredentialsApi,
    Configuration,
    CreateGithubCredentialsParams
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let body: CreateGithubCredentialsParams; //Parameters used when creating a GitHub credential.

const { status, data } = await apiInstance.createCredentials(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateGithubCredentialsParams**| Parameters used when creating a GitHub credential. | |


### Return type

**ForgeCredentials**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeCredentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createGiteaCredentials**
> ForgeCredentials createGiteaCredentials(body)


### Example

```typescript
import {
    CredentialsApi,
    Configuration,
    CreateGiteaCredentialsParams
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let body: CreateGiteaCredentialsParams; //Parameters used when creating a Gitea credential.

const { status, data } = await apiInstance.createGiteaCredentials(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateGiteaCredentialsParams**| Parameters used when creating a Gitea credential. | |


### Return type

**ForgeCredentials**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeCredentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCredentials**
> APIErrorResponse deleteCredentials()


### Example

```typescript
import {
    CredentialsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let id: number; //ID of the GitHub credential. (default to undefined)

const { status, data } = await apiInstance.deleteCredentials(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | ID of the GitHub credential. | defaults to undefined|


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

# **deleteGiteaCredentials**
> APIErrorResponse deleteGiteaCredentials()


### Example

```typescript
import {
    CredentialsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let id: number; //ID of the Gitea credential. (default to undefined)

const { status, data } = await apiInstance.deleteGiteaCredentials(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | ID of the Gitea credential. | defaults to undefined|


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

# **getCredentials**
> ForgeCredentials getCredentials()


### Example

```typescript
import {
    CredentialsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let id: number; //ID of the GitHub credential. (default to undefined)

const { status, data } = await apiInstance.getCredentials(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | ID of the GitHub credential. | defaults to undefined|


### Return type

**ForgeCredentials**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeCredentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getGiteaCredentials**
> ForgeCredentials getGiteaCredentials()


### Example

```typescript
import {
    CredentialsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let id: number; //ID of the Gitea credential. (default to undefined)

const { status, data } = await apiInstance.getGiteaCredentials(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | ID of the Gitea credential. | defaults to undefined|


### Return type

**ForgeCredentials**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeCredentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listCredentials**
> Array<ForgeCredentials> listCredentials()


### Example

```typescript
import {
    CredentialsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

const { status, data } = await apiInstance.listCredentials();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ForgeCredentials>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Credentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listGiteaCredentials**
> Array<ForgeCredentials> listGiteaCredentials()


### Example

```typescript
import {
    CredentialsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

const { status, data } = await apiInstance.listGiteaCredentials();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ForgeCredentials>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Credentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateCredentials**
> ForgeCredentials updateCredentials(body)


### Example

```typescript
import {
    CredentialsApi,
    Configuration,
    UpdateGithubCredentialsParams
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let id: number; //ID of the GitHub credential. (default to undefined)
let body: UpdateGithubCredentialsParams; //Parameters used when updating a GitHub credential.

const { status, data } = await apiInstance.updateCredentials(
    id,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateGithubCredentialsParams**| Parameters used when updating a GitHub credential. | |
| **id** | [**number**] | ID of the GitHub credential. | defaults to undefined|


### Return type

**ForgeCredentials**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeCredentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGiteaCredentials**
> ForgeCredentials updateGiteaCredentials(body)


### Example

```typescript
import {
    CredentialsApi,
    Configuration,
    UpdateGiteaCredentialsParams
} from './api';

const configuration = new Configuration();
const apiInstance = new CredentialsApi(configuration);

let id: number; //ID of the Gitea credential. (default to undefined)
let body: UpdateGiteaCredentialsParams; //Parameters used when updating a Gitea credential.

const { status, data } = await apiInstance.updateGiteaCredentials(
    id,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateGiteaCredentialsParams**| Parameters used when updating a Gitea credential. | |
| **id** | [**number**] | ID of the Gitea credential. | defaults to undefined|


### Return type

**ForgeCredentials**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeCredentials |  -  |
|**400** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

