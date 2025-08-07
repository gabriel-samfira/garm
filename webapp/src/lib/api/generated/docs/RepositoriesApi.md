# RepositoriesApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createRepo**](#createrepo) | **POST** /repositories | Create repository with the parameters given.|
|[**createRepoPool**](#createrepopool) | **POST** /repositories/{repoID}/pools | Create repository pool with the parameters given.|
|[**createRepoScaleSet**](#createreposcaleset) | **POST** /repositories/{repoID}/scalesets | Create repository scale set with the parameters given.|
|[**deleteRepo**](#deleterepo) | **DELETE** /repositories/{repoID} | Delete repository by ID.|
|[**deleteRepoPool**](#deleterepopool) | **DELETE** /repositories/{repoID}/pools/{poolID} | Delete repository pool by ID.|
|[**getRepo**](#getrepo) | **GET** /repositories/{repoID} | Get repository by ID.|
|[**getRepoPool**](#getrepopool) | **GET** /repositories/{repoID}/pools/{poolID} | Get repository pool by ID.|
|[**getRepoWebhookInfo**](#getrepowebhookinfo) | **GET** /repositories/{repoID}/webhook | Get information about the GARM installed webhook on a repository.|
|[**installRepoWebhook**](#installrepowebhook) | **POST** /repositories/{repoID}/webhook | |
|[**listRepoInstances**](#listrepoinstances) | **GET** /repositories/{repoID}/instances | List repository instances.|
|[**listRepoPools**](#listrepopools) | **GET** /repositories/{repoID}/pools | List repository pools.|
|[**listRepoScaleSets**](#listreposcalesets) | **GET** /repositories/{repoID}/scalesets | List repository scale sets.|
|[**listRepos**](#listrepos) | **GET** /repositories | List repositories.|
|[**uninstallRepoWebhook**](#uninstallrepowebhook) | **DELETE** /repositories/{repoID}/webhook | Uninstall organization webhook.|
|[**updateRepo**](#updaterepo) | **PUT** /repositories/{repoID} | Update repository with the parameters given.|
|[**updateRepoPool**](#updaterepopool) | **PUT** /repositories/{repoID}/pools/{poolID} | Update repository pool with the parameters given.|

# **createRepo**
> Repository createRepo(body)


### Example

```typescript
import {
    RepositoriesApi,
    Configuration,
    CreateRepoParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let body: CreateRepoParams; //Parameters used when creating the repository.

const { status, data } = await apiInstance.createRepo(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateRepoParams**| Parameters used when creating the repository. | |


### Return type

**Repository**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Repository |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createRepoPool**
> Pool createRepoPool(body)


### Example

```typescript
import {
    RepositoriesApi,
    Configuration,
    CreatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)
let body: CreatePoolParams; //Parameters used when creating the repository pool.

const { status, data } = await apiInstance.createRepoPool(
    repoID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreatePoolParams**| Parameters used when creating the repository pool. | |
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**Pool**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pool |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createRepoScaleSet**
> ScaleSet createRepoScaleSet(body)


### Example

```typescript
import {
    RepositoriesApi,
    Configuration,
    CreateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)
let body: CreateScaleSetParams; //Parameters used when creating the repository scale set.

const { status, data } = await apiInstance.createRepoScaleSet(
    repoID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateScaleSetParams**| Parameters used when creating the repository scale set. | |
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**ScaleSet**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ScaleSet |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRepo**
> APIErrorResponse deleteRepo()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //ID of the repository to delete. (default to undefined)
let keepWebhook: boolean; //If true and a webhook is installed for this repo, it will not be removed. (optional) (default to undefined)

const { status, data } = await apiInstance.deleteRepo(
    repoID,
    keepWebhook
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | ID of the repository to delete. | defaults to undefined|
| **keepWebhook** | [**boolean**] | If true and a webhook is installed for this repo, it will not be removed. | (optional) defaults to undefined|


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

# **deleteRepoPool**
> APIErrorResponse deleteRepoPool()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)
let poolID: string; //ID of the repository pool to delete. (default to undefined)

const { status, data } = await apiInstance.deleteRepoPool(
    repoID,
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|
| **poolID** | [**string**] | ID of the repository pool to delete. | defaults to undefined|


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

# **getRepo**
> Repository getRepo()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //ID of the repository to fetch. (default to undefined)

const { status, data } = await apiInstance.getRepo(
    repoID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | ID of the repository to fetch. | defaults to undefined|


### Return type

**Repository**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Repository |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRepoPool**
> Pool getRepoPool()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)
let poolID: string; //Pool ID. (default to undefined)

const { status, data } = await apiInstance.getRepoPool(
    repoID,
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|
| **poolID** | [**string**] | Pool ID. | defaults to undefined|


### Return type

**Pool**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pool |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRepoWebhookInfo**
> HookInfo getRepoWebhookInfo()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

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

# **installRepoWebhook**
> HookInfo installRepoWebhook(body)

Install the GARM webhook for an organization. The secret configured on the organization will be used to validate the requests.

### Example

```typescript
import {
    RepositoriesApi,
    Configuration,
    InstallWebhookParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

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

# **listRepoInstances**
> Array<Instance> listRepoInstances()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)

const { status, data } = await apiInstance.listRepoInstances(
    repoID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**Array<Instance>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Instances |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listRepoPools**
> Array<Pool> listRepoPools()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)

const { status, data } = await apiInstance.listRepoPools(
    repoID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**Array<Pool>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pools |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listRepoScaleSets**
> Array<ScaleSet> listRepoScaleSets()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)

const { status, data } = await apiInstance.listRepoScaleSets(
    repoID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **repoID** | [**string**] | Repository ID. | defaults to undefined|


### Return type

**Array<ScaleSet>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ScaleSets |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listRepos**
> Array<Repository> listRepos()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let owner: string; //Exact owner name to filter by (optional) (default to undefined)
let name: string; //Exact repository name to filter by (optional) (default to undefined)
let endpoint: string; //Exact endpoint name to filter by (optional) (default to undefined)

const { status, data } = await apiInstance.listRepos(
    owner,
    name,
    endpoint
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **owner** | [**string**] | Exact owner name to filter by | (optional) defaults to undefined|
| **name** | [**string**] | Exact repository name to filter by | (optional) defaults to undefined|
| **endpoint** | [**string**] | Exact endpoint name to filter by | (optional) defaults to undefined|


### Return type

**Array<Repository>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Repositories |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uninstallRepoWebhook**
> APIErrorResponse uninstallRepoWebhook()


### Example

```typescript
import {
    RepositoriesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

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

# **updateRepo**
> Repository updateRepo(body)


### Example

```typescript
import {
    RepositoriesApi,
    Configuration,
    UpdateEntityParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //ID of the repository to update. (default to undefined)
let body: UpdateEntityParams; //Parameters used when updating the repository.

const { status, data } = await apiInstance.updateRepo(
    repoID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateEntityParams**| Parameters used when updating the repository. | |
| **repoID** | [**string**] | ID of the repository to update. | defaults to undefined|


### Return type

**Repository**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Repository |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRepoPool**
> Pool updateRepoPool(body)


### Example

```typescript
import {
    RepositoriesApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RepositoriesApi(configuration);

let repoID: string; //Repository ID. (default to undefined)
let poolID: string; //ID of the repository pool to update. (default to undefined)
let body: UpdatePoolParams; //Parameters used when updating the repository pool.

const { status, data } = await apiInstance.updateRepoPool(
    repoID,
    poolID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdatePoolParams**| Parameters used when updating the repository pool. | |
| **repoID** | [**string**] | Repository ID. | defaults to undefined|
| **poolID** | [**string**] | ID of the repository pool to update. | defaults to undefined|


### Return type

**Pool**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pool |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

