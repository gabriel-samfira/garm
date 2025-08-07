# InstancesApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**deleteInstance**](#deleteinstance) | **DELETE** /instances/{instanceName} | Delete runner instance by name.|
|[**getInstance**](#getinstance) | **GET** /instances/{instanceName} | Get runner instance by name.|
|[**listEnterpriseInstances**](#listenterpriseinstances) | **GET** /enterprises/{enterpriseID}/instances | List enterprise instances.|
|[**listInstances**](#listinstances) | **GET** /instances | Get all runners\&#39; instances.|
|[**listOrgInstances**](#listorginstances) | **GET** /organizations/{orgID}/instances | List organization instances.|
|[**listPoolInstances**](#listpoolinstances) | **GET** /pools/{poolID}/instances | List runner instances in a pool.|
|[**listRepoInstances**](#listrepoinstances) | **GET** /repositories/{repoID}/instances | List repository instances.|
|[**listScaleSetInstances**](#listscalesetinstances) | **GET** /scalesets/{scalesetID}/instances | List runner instances in a scale set.|

# **deleteInstance**
> APIErrorResponse deleteInstance()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

let instanceName: string; //Runner instance name. (default to undefined)
let forceRemove: boolean; //If true GARM will ignore any provider error when removing the runner and will continue to remove the runner from github and the GARM database. (optional) (default to undefined)
let bypassGHUnauthorized: boolean; //If true GARM will ignore unauthorized errors returned by GitHub when removing a runner. This is useful if you want to clean up runners and your credentials have expired. (optional) (default to undefined)

const { status, data } = await apiInstance.deleteInstance(
    instanceName,
    forceRemove,
    bypassGHUnauthorized
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **instanceName** | [**string**] | Runner instance name. | defaults to undefined|
| **forceRemove** | [**boolean**] | If true GARM will ignore any provider error when removing the runner and will continue to remove the runner from github and the GARM database. | (optional) defaults to undefined|
| **bypassGHUnauthorized** | [**boolean**] | If true GARM will ignore unauthorized errors returned by GitHub when removing a runner. This is useful if you want to clean up runners and your credentials have expired. | (optional) defaults to undefined|


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

# **getInstance**
> Instance getInstance()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

let instanceName: string; //Runner instance name. (default to undefined)

const { status, data } = await apiInstance.getInstance(
    instanceName
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **instanceName** | [**string**] | Runner instance name. | defaults to undefined|


### Return type

**Instance**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Instance |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listEnterpriseInstances**
> Array<Instance> listEnterpriseInstances()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

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

# **listInstances**
> Array<Instance> listInstances()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

const { status, data } = await apiInstance.listInstances();
```

### Parameters
This endpoint does not have any parameters.


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

# **listOrgInstances**
> Array<Instance> listOrgInstances()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

let orgID: string; //Organization ID. (default to undefined)

const { status, data } = await apiInstance.listOrgInstances(
    orgID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | Organization ID. | defaults to undefined|


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

# **listPoolInstances**
> Array<Instance> listPoolInstances()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

let poolID: string; //Runner pool ID. (default to undefined)

const { status, data } = await apiInstance.listPoolInstances(
    poolID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **poolID** | [**string**] | Runner pool ID. | defaults to undefined|


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

# **listRepoInstances**
> Array<Instance> listRepoInstances()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

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

# **listScaleSetInstances**
> Array<Instance> listScaleSetInstances()


### Example

```typescript
import {
    InstancesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new InstancesApi(configuration);

let scalesetID: string; //Runner scale set ID. (default to undefined)

const { status, data } = await apiInstance.listScaleSetInstances(
    scalesetID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **scalesetID** | [**string**] | Runner scale set ID. | defaults to undefined|


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

