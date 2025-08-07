# ControllerInfo


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**callback_url** | **string** | CallbackURL is the URL where instances can send updates back to the controller. This URL is used by instances to send status updates back to the controller. The URL itself may be made available to instances via a reverse proxy or a load balancer. That means that the user is responsible for telling GARM what the public URL is, by setting this field. | [optional] [default to undefined]
**controller_id** | **string** | ControllerID is the unique ID of this controller. This ID gets generated automatically on controller init. | [optional] [default to undefined]
**controller_webhook_url** | **string** | ControllerWebhookURL is the controller specific URL where webhooks will be received. This field holds the WebhookURL defined above to which we append the ControllerID. Functionally it is the same as WebhookURL, but it allows us to safely manage webhooks from GARM without accidentally removing webhooks from other services or GARM controllers. | [optional] [default to undefined]
**hostname** | **string** | Hostname is the hostname of the machine that runs this controller. In the future, this field will be migrated to a separate table that will keep track of each the controller nodes that are part of a cluster. This will happen when we implement controller scale-out capability. | [optional] [default to undefined]
**metadata_url** | **string** | MetadataURL is the public metadata URL of the GARM instance. This URL is used by instances to fetch information they need to set themselves up. The URL itself may be made available to runners via a reverse proxy or a load balancer. That means that the user is responsible for telling GARM what the public URL is, by setting this field. | [optional] [default to undefined]
**minimum_job_age_backoff** | **number** | MinimumJobAgeBackoff is the minimum time in seconds that a job must be in queued state before GARM will attempt to allocate a runner for it. When set to a non zero value, GARM will ignore the job until the job\&#39;s age is greater than this value. When using the min_idle_runners feature of a pool, this gives enough time for potential idle runners to pick up the job before GARM attempts to allocate a new runner, thus avoiding the need to potentially scale down runners later. | [optional] [default to undefined]
**version** | **string** | Version is the version of the GARM controller. | [optional] [default to undefined]
**webhook_url** | **string** | WebhookURL is the base URL where the controller will receive webhooks from github. When webhook management is used, this URL is used as a base to which the controller UUID is appended and which will receive the webhooks. The URL itself may be made available to instances via a reverse proxy or a load balancer. That means that the user is responsible for telling GARM what the public URL is, by setting this field. | [optional] [default to undefined]

## Example

```typescript
import { ControllerInfo } from './api';

const instance: ControllerInfo = {
    callback_url,
    controller_id,
    controller_webhook_url,
    hostname,
    metadata_url,
    minimum_job_age_backoff,
    version,
    webhook_url,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
