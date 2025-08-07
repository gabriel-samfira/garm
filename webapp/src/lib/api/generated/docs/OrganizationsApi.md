# OrganizationsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createOrg**](#createorg) | **POST** /organizations | Create organization with the parameters given.|
|[**createOrgPool**](#createorgpool) | **POST** /organizations/{orgID}/pools | Create organization pool with the parameters given.|
|[**createOrgScaleSet**](#createorgscaleset) | **POST** /organizations/{orgID}/scalesets | Create organization scale set with the parameters given.|
|[**deleteOrg**](#deleteorg) | **DELETE** /organizations/{orgID} | Delete organization by ID.|
|[**deleteOrgPool**](#deleteorgpool) | **DELETE** /organizations/{orgID}/pools/{poolID} | Delete organization pool by ID.|
|[**getOrg**](#getorg) | **GET** /organizations/{orgID} | Get organization by ID.|
|[**getOrgPool**](#getorgpool) | **GET** /organizations/{orgID}/pools/{poolID} | Get organization pool by ID.|
|[**getOrgWebhookInfo**](#getorgwebhookinfo) | **GET** /organizations/{orgID}/webhook | Get information about the GARM installed webhook on an organization.|
|[**installOrgWebhook**](#installorgwebhook) | **POST** /organizations/{orgID}/webhook | |
|[**listOrgInstances**](#listorginstances) | **GET** /organizations/{orgID}/instances | List organization instances.|
|[**listOrgPools**](#listorgpools) | **GET** /organizations/{orgID}/pools | List organization pools.|
|[**listOrgScaleSets**](#listorgscalesets) | **GET** /organizations/{orgID}/scalesets | List organization scale sets.|
|[**listOrgs**](#listorgs) | **GET** /organizations | List organizations.|
|[**uninstallOrgWebhook**](#uninstallorgwebhook) | **DELETE** /organizations/{orgID}/webhook | Uninstall organization webhook.|
|[**updateOrg**](#updateorg) | **PUT** /organizations/{orgID} | Update organization with the parameters given.|
|[**updateOrgPool**](#updateorgpool) | **PUT** /organizations/{orgID}/pools/{poolID} | Update organization pool with the parameters given.|

# **createOrg**
> Organization createOrg(body)


### Example

```typescript
import {
    OrganizationsApi,
    Configuration,
    CreateOrgParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

let body: CreateOrgParams; //Parameters used when creating the organization.

const { status, data } = await apiInstance.createOrg(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **CreateOrgParams**| Parameters used when creating the organization. | |


### Return type

**Organization**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Organization |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createOrgPool**
> Pool createOrgPool(body)


### Example

```typescript
import {
    OrganizationsApi,
    Configuration,
    CreatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **createOrgScaleSet**
> ScaleSet createOrgScaleSet(body)


### Example

```typescript
import {
    OrganizationsApi,
    Configuration,
    CreateScaleSetParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **deleteOrg**
> APIErrorResponse deleteOrg()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

let orgID: string; //ID of the organization to delete. (default to undefined)
let keepWebhook: boolean; //If true and a webhook is installed for this organization, it will not be removed. (optional) (default to undefined)

const { status, data } = await apiInstance.deleteOrg(
    orgID,
    keepWebhook
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | ID of the organization to delete. | defaults to undefined|
| **keepWebhook** | [**boolean**] | If true and a webhook is installed for this organization, it will not be removed. | (optional) defaults to undefined|


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
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **getOrg**
> Organization getOrg()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

let orgID: string; //ID of the organization to fetch. (default to undefined)

const { status, data } = await apiInstance.getOrg(
    orgID
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orgID** | [**string**] | ID of the organization to fetch. | defaults to undefined|


### Return type

**Organization**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Organization |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getOrgPool**
> Pool getOrgPool()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **getOrgWebhookInfo**
> HookInfo getOrgWebhookInfo()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **installOrgWebhook**
> HookInfo installOrgWebhook(body)

Install the GARM webhook for an organization. The secret configured on the organization will be used to validate the requests.

### Example

```typescript
import {
    OrganizationsApi,
    Configuration,
    InstallWebhookParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **listOrgInstances**
> Array<Instance> listOrgInstances()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **listOrgPools**
> Array<Pool> listOrgPools()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **listOrgScaleSets**
> Array<ScaleSet> listOrgScaleSets()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **listOrgs**
> Array<Organization> listOrgs()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

let name: string; //Exact organization name to filter by (optional) (default to undefined)
let endpoint: string; //Exact endpoint name to filter by (optional) (default to undefined)

const { status, data } = await apiInstance.listOrgs(
    name,
    endpoint
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] | Exact organization name to filter by | (optional) defaults to undefined|
| **endpoint** | [**string**] | Exact endpoint name to filter by | (optional) defaults to undefined|


### Return type

**Array<Organization>**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Organizations |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uninstallOrgWebhook**
> APIErrorResponse uninstallOrgWebhook()


### Example

```typescript
import {
    OrganizationsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

# **updateOrg**
> Organization updateOrg(body)


### Example

```typescript
import {
    OrganizationsApi,
    Configuration,
    UpdateEntityParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

let orgID: string; //ID of the organization to update. (default to undefined)
let body: UpdateEntityParams; //Parameters used when updating the organization.

const { status, data } = await apiInstance.updateOrg(
    orgID,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **UpdateEntityParams**| Parameters used when updating the organization. | |
| **orgID** | [**string**] | ID of the organization to update. | defaults to undefined|


### Return type

**Organization**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Organization |  -  |
|**0** | APIErrorResponse |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateOrgPool**
> Pool updateOrgPool(body)


### Example

```typescript
import {
    OrganizationsApi,
    Configuration,
    UpdatePoolParams
} from './api';

const configuration = new Configuration();
const apiInstance = new OrganizationsApi(configuration);

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

