# EnterprisesApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createEnterprise**](#createenterprise) | **POST** /enterprises | Create enterprise with the given parameters.|
|[**createEnterprisePool**](#createenterprisepool) | **POST** /enterprises/{enterpriseID}/pools | Create enterprise pool with the parameters given.|
|[**createEnterpriseScaleSet**](#createenterprisescaleset) | **POST** /enterprises/{enterpriseID}/scalesets | Create enterprise pool with the parameters given.|
|[**deleteEnterprise**](#deleteenterprise) | **DELETE** /enterprises/{enterpriseID} | Delete enterprise by ID.|
|[**deleteEnterprisePool**](#deleteenterprisepool) | **DELETE** /enterprises/{enterpriseID}/pools/{poolID} | Delete enterprise pool by ID.|
|[**getEnterprise**](#getenterprise) | **GET** /enterprises/{enterpriseID} | Get enterprise by ID.|
|[**getEnterprisePool**](#getenterprisepool) | **GET** /enterprises/{enterpriseID}/pools/{poolID} | Get enterprise pool by ID.|
|[**listEnterpriseInstances**](#listenterpriseinstances) | **GET** /enterprises/{enterpriseID}/instances | List enterprise instances.|
|[**listEnterprisePools**](#listenterprisepools) | **GET** /enterprises/{enterpriseID}/pools | List enterprise pools.|
|[**listEnterpriseScaleSets**](#listenterprisescalesets) | **GET** /enterprises/{enterpriseID}/scalesets | List enterprise scale sets.|
|[**listEnterprises**](#listenterprises) | **GET** /enterprises | List all enterprises.|
|[**updateEnterprise**](#updateenterprise) | **PUT** /enterprises/{enterpriseID} | Update enterprise with the given parameters.|
|[**updateEnterprisePool**](#updateenterprisepool) | **PUT** /enterprises/{enterpriseID}/pools/{poolID} | Update enterprise pool with the parameters given.|

# **createEnterprise**
> Enterprise createEnterprise(body)


### Example

```typescript
import {
    EnterprisesApi,
    Configuration,
    CreateEnterpriseParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let body: CreateEnterpriseParams; //Parameters used to create the enterprise.

const { status, data } = await apiInstance.createEnterprise(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateEnterpriseParams**| Parameters used to create the enterprise. | |


### Return type

**Enterprise**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Enterprise |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createEnterprisePool**
> Pool createEnterprisePool(body)


### Example

```typescript
import {
    EnterprisesApi,
    Configuration,
    CreatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

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

# **createEnterpriseScaleSet**
> ScaleSet createEnterpriseScaleSet(body)


### Example

```typescript
import {
    EnterprisesApi,
    Configuration,
    CreateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)
let body: CreateScaleSetParams; //Parameters used when creating the enterprise scale set.

const { status, data } = await apiInstance.createEnterpriseScaleSet(
    enterpriseID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateScaleSetParams**| Parameters used when creating the enterprise scale set. | |
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|


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

# **deleteEnterprise**
> APIErrorResponse deleteEnterprise()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let enterpriseID: string; //ID of the enterprise to delete. (default to undefined)

const { status, data } = await apiInstance.deleteEnterprise(
    enterpriseID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | ID of the enterprise to delete. | defaults to undefined|


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

# **deleteEnterprisePool**
> APIErrorResponse deleteEnterprisePool()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

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

# **getEnterprise**
> Enterprise getEnterprise()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let enterpriseID: string; //The ID of the enterprise to fetch. (default to undefined)

const { status, data } = await apiInstance.getEnterprise(
    enterpriseID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | The ID of the enterprise to fetch. | defaults to undefined|


### Return type

**Enterprise**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Enterprise |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getEnterprisePool**
> Pool getEnterprisePool()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

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

# **listEnterpriseInstances**
> Array<Instance> listEnterpriseInstances()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)

const { status, data } = await apiInstance.listEnterpriseInstances(
    enterpriseID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|


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

# **listEnterprisePools**
> Array<Pool> listEnterprisePools()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

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

# **listEnterpriseScaleSets**
> Array<ScaleSet> listEnterpriseScaleSets()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let enterpriseID: string; //Enterprise ID. (default to undefined)

const { status, data } = await apiInstance.listEnterpriseScaleSets(
    enterpriseID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **enterpriseID** | [**string**] | Enterprise ID. | defaults to undefined|


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

# **listEnterprises**
> Array<Enterprise> listEnterprises()


### Example

```typescript
import {
    EnterprisesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let name: string; //Exact enterprise name to filter by (optional) (default to undefined)
let endpoint: string; //Exact endpoint name to filter by (optional) (default to undefined)

const { status, data } = await apiInstance.listEnterprises(
    name,
    endpoint
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | Exact enterprise name to filter by | (optional) defaults to undefined|
| **endpoint** | [**string**] | Exact endpoint name to filter by | (optional) defaults to undefined|


### Return type

**Array<Enterprise>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Enterprises |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateEnterprise**
> Enterprise updateEnterprise(body)


### Example

```typescript
import {
    EnterprisesApi,
    Configuration,
    UpdateEntityParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

let enterpriseID: string; //The ID of the enterprise to update. (default to undefined)
let body: UpdateEntityParams; //Parameters used when updating the enterprise.

const { status, data } = await apiInstance.updateEnterprise(
    enterpriseID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateEntityParams**| Parameters used when updating the enterprise. | |
| **enterpriseID** | [**string**] | The ID of the enterprise to update. | defaults to undefined|


### Return type

**Enterprise**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Enterprise |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateEnterprisePool**
> Pool updateEnterprisePool(body)


### Example

```typescript
import {
    EnterprisesApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new EnterprisesApi(configuration);

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

