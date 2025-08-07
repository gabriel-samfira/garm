# EndpointsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createGiteaEndpoint**](#creategiteaendpoint) | **POST** /gitea/endpoints | Create a Gitea Endpoint.|
|[**createGithubEndpoint**](#creategithubendpoint) | **POST** /github/endpoints | Create a GitHub Endpoint.|
|[**deleteGiteaEndpoint**](#deletegiteaendpoint) | **DELETE** /gitea/endpoints/{name} | Delete a Gitea Endpoint.|
|[**deleteGithubEndpoint**](#deletegithubendpoint) | **DELETE** /github/endpoints/{name} | Delete a GitHub Endpoint.|
|[**getGiteaEndpoint**](#getgiteaendpoint) | **GET** /gitea/endpoints/{name} | Get a Gitea Endpoint.|
|[**getGithubEndpoint**](#getgithubendpoint) | **GET** /github/endpoints/{name} | Get a GitHub Endpoint.|
|[**listGiteaEndpoints**](#listgiteaendpoints) | **GET** /gitea/endpoints | List all Gitea Endpoints.|
|[**listGithubEndpoints**](#listgithubendpoints) | **GET** /github/endpoints | List all GitHub Endpoints.|
|[**updateGiteaEndpoint**](#updategiteaendpoint) | **PUT** /gitea/endpoints/{name} | Update a Gitea Endpoint.|
|[**updateGithubEndpoint**](#updategithubendpoint) | **PUT** /github/endpoints/{name} | Update a GitHub Endpoint.|

# **createGiteaEndpoint**
> ForgeEndpoint createGiteaEndpoint(body)


### Example

```typescript
import {
    EndpointsApi,
    Configuration,
    CreateGiteaEndpointParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let body: CreateGiteaEndpointParams; //Parameters used when creating a Gitea endpoint.

const { status, data } = await apiInstance.createGiteaEndpoint(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateGiteaEndpointParams**| Parameters used when creating a Gitea endpoint. | |


### Return type

**ForgeEndpoint**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoint |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createGithubEndpoint**
> ForgeEndpoint createGithubEndpoint(body)


### Example

```typescript
import {
    EndpointsApi,
    Configuration,
    CreateGithubEndpointParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let body: CreateGithubEndpointParams; //Parameters used when creating a GitHub endpoint.

const { status, data } = await apiInstance.createGithubEndpoint(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateGithubEndpointParams**| Parameters used when creating a GitHub endpoint. | |


### Return type

**ForgeEndpoint**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoint |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteGiteaEndpoint**
> APIErrorResponse deleteGiteaEndpoint()


### Example

```typescript
import {
    EndpointsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let name: string; //The name of the Gitea endpoint. (default to undefined)

const { status, data } = await apiInstance.deleteGiteaEndpoint(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | The name of the Gitea endpoint. | defaults to undefined|


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

# **deleteGithubEndpoint**
> APIErrorResponse deleteGithubEndpoint()


### Example

```typescript
import {
    EndpointsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let name: string; //The name of the GitHub endpoint. (default to undefined)

const { status, data } = await apiInstance.deleteGithubEndpoint(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | The name of the GitHub endpoint. | defaults to undefined|


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

# **getGiteaEndpoint**
> ForgeEndpoint getGiteaEndpoint()


### Example

```typescript
import {
    EndpointsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let name: string; //The name of the Gitea endpoint. (default to undefined)

const { status, data } = await apiInstance.getGiteaEndpoint(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | The name of the Gitea endpoint. | defaults to undefined|


### Return type

**ForgeEndpoint**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoint |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getGithubEndpoint**
> ForgeEndpoint getGithubEndpoint()


### Example

```typescript
import {
    EndpointsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let name: string; //The name of the GitHub endpoint. (default to undefined)

const { status, data } = await apiInstance.getGithubEndpoint(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | The name of the GitHub endpoint. | defaults to undefined|


### Return type

**ForgeEndpoint**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoint |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listGiteaEndpoints**
> Array<ForgeEndpoint> listGiteaEndpoints()


### Example

```typescript
import {
    EndpointsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

const { status, data } = await apiInstance.listGiteaEndpoints();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ForgeEndpoint>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoints |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listGithubEndpoints**
> Array<ForgeEndpoint> listGithubEndpoints()


### Example

```typescript
import {
    EndpointsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

const { status, data } = await apiInstance.listGithubEndpoints();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ForgeEndpoint>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoints |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGiteaEndpoint**
> ForgeEndpoint updateGiteaEndpoint(body)


### Example

```typescript
import {
    EndpointsApi,
    Configuration,
    UpdateGiteaEndpointParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let name: string; //The name of the Gitea endpoint. (default to undefined)
let body: UpdateGiteaEndpointParams; //Parameters used when updating a Gitea endpoint.

const { status, data } = await apiInstance.updateGiteaEndpoint(
    name,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateGiteaEndpointParams**| Parameters used when updating a Gitea endpoint. | |
| **name** | [**string**] | The name of the Gitea endpoint. | defaults to undefined|


### Return type

**ForgeEndpoint**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoint |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGithubEndpoint**
> ForgeEndpoint updateGithubEndpoint(body)


### Example

```typescript
import {
    EndpointsApi,
    Configuration,
    UpdateGithubEndpointParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EndpointsApi(configuration);

let name: string; //The name of the GitHub endpoint. (default to undefined)
let body: UpdateGithubEndpointParams; //Parameters used when updating a GitHub endpoint.

const { status, data } = await apiInstance.updateGithubEndpoint(
    name,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateGithubEndpointParams**| Parameters used when updating a GitHub endpoint. | |
| **name** | [**string**] | The name of the GitHub endpoint. | defaults to undefined|


### Return type

**ForgeEndpoint**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ForgeEndpoint |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

