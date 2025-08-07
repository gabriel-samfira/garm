# PoolsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createEnterprisePool**](#createenterprisepool) | **POST** /enterprises/{enterpriseID}/pools | Create enterprise pool with the parameters given.|
|[**createOrgPool**](#createorgpool) | **POST** /organizations/{orgID}/pools | Create organization pool with the parameters given.|
|[**createRepoPool**](#createrepopool) | **POST** /repositories/{repoID}/pools | Create repository pool with the parameters given.|
|[**deleteEnterprisePool**](#deleteenterprisepool) | **DELETE** /enterprises/{enterpriseID}/pools/{poolID} | Delete enterprise pool by ID.|
|[**deleteOrgPool**](#deleteorgpool) | **DELETE** /organizations/{orgID}/pools/{poolID} | Delete organization pool by ID.|
|[**deletePool**](#deletepool) | **DELETE** /pools/{poolID} | Delete pool by ID.|
|[**deleteRepoPool**](#deleterepopool) | **DELETE** /repositories/{repoID}/pools/{poolID} | Delete repository pool by ID.|
|[**getEnterprisePool**](#getenterprisepool) | **GET** /enterprises/{enterpriseID}/pools/{poolID} | Get enterprise pool by ID.|
|[**getOrgPool**](#getorgpool) | **GET** /organizations/{orgID}/pools/{poolID} | Get organization pool by ID.|
|[**getPool**](#getpool) | **GET** /pools/{poolID} | Get pool by ID.|
|[**getRepoPool**](#getrepopool) | **GET** /repositories/{repoID}/pools/{poolID} | Get repository pool by ID.|
|[**listEnterprisePools**](#listenterprisepools) | **GET** /enterprises/{enterpriseID}/pools | List enterprise pools.|
|[**listOrgPools**](#listorgpools) | **GET** /organizations/{orgID}/pools | List organization pools.|
|[**listPools**](#listpools) | **GET** /pools | List all pools.|
|[**listRepoPools**](#listrepopools) | **GET** /repositories/{repoID}/pools | List repository pools.|
|[**updateEnterprisePool**](#updateenterprisepool) | **PUT** /enterprises/{enterpriseID}/pools/{poolID} | Update enterprise pool with the parameters given.|
|[**updateOrgPool**](#updateorgpool) | **PUT** /organizations/{orgID}/pools/{poolID} | Update organization pool with the parameters given.|
|[**updatePool**](#updatepool) | **PUT** /pools/{poolID} | Update pool by ID.|
|[**updateRepoPool**](#updaterepopool) | **PUT** /repositories/{repoID}/pools/{poolID} | Update repository pool with the parameters given.|

# **createEnterprisePool**
> Pool createEnterprisePool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    CreatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)
let body: CreatePoolParams; //Parameters used when creating the enterprise pool.

const { status, data } = await apiInstance.createEnterprisePool(
    enterpriseID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreatePoolParams**| Parameters used when creating the enterprise pool. | |
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|


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

# **createOrgPool**
> Pool createOrgPool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    CreatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)
let body: CreatePoolParams; //Parameters used when creating the organization pool.

const { status, data } = await apiInstance.createOrgPool(
    orgID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreatePoolParams**| Parameters used when creating the organization pool. | |
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


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

# **createRepoPool**
> Pool createRepoPool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    CreatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

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

# **deleteEnterprisePool**
> APIErrorResponse deleteEnterprisePool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)
let poolID: string; //ID of the enterprise pool to delete. (default to undefined)

const { status, data } = await apiInstance.deleteEnterprisePool(
    enterpriseID,
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|
| **poolID** | [**string**] | ID of the enterprise pool to delete. | defaults to undefined|


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

# **deleteOrgPool**
> APIErrorResponse deleteOrgPool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)
let poolID: string; //ID of the organization pool to delete. (default to undefined)

const { status, data } = await apiInstance.deleteOrgPool(
    orgID,
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|
| **poolID** | [**string**] | ID of the organization pool to delete. | defaults to undefined|


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

# **deletePool**
> APIErrorResponse deletePool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let poolID: string; //ID of the pool to delete. (default to undefined)

const { status, data } = await apiInstance.deletePool(
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **poolID** | [**string**] | ID of the pool to delete. | defaults to undefined|


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
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

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

# **getEnterprisePool**
> Pool getEnterprisePool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)
let poolID: string; //Pool ID. (default to undefined)

const { status, data } = await apiInstance.getEnterprisePool(
    enterpriseID,
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|
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

# **getOrgPool**
> Pool getOrgPool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)
let poolID: string; //Pool ID. (default to undefined)

const { status, data } = await apiInstance.getOrgPool(
    orgID,
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|
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

# **getPool**
> Pool getPool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let poolID: string; //ID of the pool to fetch. (default to undefined)

const { status, data } = await apiInstance.getPool(
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **poolID** | [**string**] | ID of the pool to fetch. | defaults to undefined|


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

# **getRepoPool**
> Pool getRepoPool()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

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

# **listEnterprisePools**
> Array<Pool> listEnterprisePools()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)

const { status, data } = await apiInstance.listEnterprisePools(
    enterpriseID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|


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

# **listOrgPools**
> Array<Pool> listOrgPools()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)

const { status, data } = await apiInstance.listOrgPools(
    orgID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


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

# **listPools**
> Array<Pool> listPools()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

const { status, data } = await apiInstance.listPools();
```

### Parameters
This endpoint does not have any parameters.


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

# **listRepoPools**
> Array<Pool> listRepoPools()


### Example

```typescript
import {
    PoolsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

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

# **updateEnterprisePool**
> Pool updateEnterprisePool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)
let poolID: string; //ID of the enterprise pool to update. (default to undefined)
let body: UpdatePoolParams; //Parameters used when updating the enterprise pool.

const { status, data } = await apiInstance.updateEnterprisePool(
    enterpriseID,
    poolID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdatePoolParams**| Parameters used when updating the enterprise pool. | |
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|
| **poolID** | [**string**] | ID of the enterprise pool to update. | defaults to undefined|


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

# **updateOrgPool**
> Pool updateOrgPool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)
let poolID: string; //ID of the organization pool to update. (default to undefined)
let body: UpdatePoolParams; //Parameters used when updating the organization pool.

const { status, data } = await apiInstance.updateOrgPool(
    orgID,
    poolID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdatePoolParams**| Parameters used when updating the organization pool. | |
| **orgID** | [**string**] | Organization ID. | defaults to undefined|
| **poolID** | [**string**] | ID of the organization pool to update. | defaults to undefined|


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

# **updatePool**
> Pool updatePool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

let poolID: string; //ID of the pool to update. (default to undefined)
let body: UpdatePoolParams; //Parameters to update the pool with.

const { status, data } = await apiInstance.updatePool(
    poolID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdatePoolParams**| Parameters to update the pool with. | |
| **poolID** | [**string**] | ID of the pool to update. | defaults to undefined|


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

# **updateRepoPool**
> Pool updateRepoPool(body)


### Example

```typescript
import {
    PoolsApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new PoolsApi(configuration);

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

