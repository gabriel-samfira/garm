# ScalesetsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createEnterpriseScaleSet**](#createenterprisescaleset) | **POST** /enterprises/{enterpriseID}/scalesets | Create enterprise pool with the parameters given.|
|[**createOrgScaleSet**](#createorgscaleset) | **POST** /organizations/{orgID}/scalesets | Create organization scale set with the parameters given.|
|[**createRepoScaleSet**](#createreposcaleset) | **POST** /repositories/{repoID}/scalesets | Create repository scale set with the parameters given.|
|[**deleteScaleSet**](#deletescaleset) | **DELETE** /scalesets/{scalesetID} | Delete scale set by ID.|
|[**getScaleSet**](#getscaleset) | **GET** /scalesets/{scalesetID} | Get scale set by ID.|
|[**listEnterpriseScaleSets**](#listenterprisescalesets) | **GET** /enterprises/{enterpriseID}/scalesets | List enterprise scale sets.|
|[**listOrgScaleSets**](#listorgscalesets) | **GET** /organizations/{orgID}/scalesets | List organization scale sets.|
|[**listRepoScaleSets**](#listreposcalesets) | **GET** /repositories/{repoID}/scalesets | List repository scale sets.|
|[**listScalesets**](#listscalesets) | **GET** /scalesets | List all scalesets.|
|[**updateScaleSet**](#updatescaleset) | **PUT** /scalesets/{scalesetID} | Update scale set by ID.|

# **createEnterpriseScaleSet**
> ScaleSet createEnterpriseScaleSet(body)


### Example

```typescript
import {
    ScalesetsApi,
    Configuration,
    CreateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

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

# **createOrgScaleSet**
> ScaleSet createOrgScaleSet(body)


### Example

```typescript
import {
    ScalesetsApi,
    Configuration,
    CreateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)
let body: CreateScaleSetParams; //Parameters used when creating the organization scale set.

const { status, data } = await apiInstance.createOrgScaleSet(
    orgID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateScaleSetParams**| Parameters used when creating the organization scale set. | |
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


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

# **createRepoScaleSet**
> ScaleSet createRepoScaleSet(body)


### Example

```typescript
import {
    ScalesetsApi,
    Configuration,
    CreateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

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

# **deleteScaleSet**
> APIErrorResponse deleteScaleSet()


### Example

```typescript
import {
    ScalesetsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

let scalesetID: string; //ID of the scale set to delete. (default to undefined)

const { status, data } = await apiInstance.deleteScaleSet(
    scalesetID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **scalesetID** | [**string**] | ID of the scale set to delete. | defaults to undefined|


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

# **getScaleSet**
> ScaleSet getScaleSet()


### Example

```typescript
import {
    ScalesetsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

let scalesetID: string; //ID of the scale set to fetch. (default to undefined)

const { status, data } = await apiInstance.getScaleSet(
    scalesetID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **scalesetID** | [**string**] | ID of the scale set to fetch. | defaults to undefined|


### Return type

**ScaleSet**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | ScaleSet |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listEnterpriseScaleSets**
> Array<ScaleSet> listEnterpriseScaleSets()


### Example

```typescript
import {
    ScalesetsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

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

# **listOrgScaleSets**
> Array<ScaleSet> listOrgScaleSets()


### Example

```typescript
import {
    ScalesetsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

let orgID: string; //Organization ID. (default to undefined)

const { status, data } = await apiInstance.listOrgScaleSets(
    orgID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


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

# **listRepoScaleSets**
> Array<ScaleSet> listRepoScaleSets()


### Example

```typescript
import {
    ScalesetsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

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

# **listScalesets**
> Array<ScaleSet> listScalesets()


### Example

```typescript
import {
    ScalesetsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

const { status, data } = await apiInstance.listScalesets();
```

### Parameters
This endpoint does not have any parameters.


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

# **updateScaleSet**
> ScaleSet updateScaleSet(body)


### Example

```typescript
import {
    ScalesetsApi,
    Configuration,
    UpdateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new ScalesetsApi(configuration);

let scalesetID: string; //ID of the scale set to update. (default to undefined)
let body: UpdateScaleSetParams; //Parameters to update the scale set with.

const { status, data } = await apiInstance.updateScaleSet(
    scalesetID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateScaleSetParams**| Parameters to update the scale set with. | |
| **scalesetID** | [**string**] | ID of the scale set to update. | defaults to undefined|


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

