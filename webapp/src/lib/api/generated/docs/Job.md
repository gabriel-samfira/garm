# Job


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**action** | **string** | Action is the specific activity that triggered the event. | [optional] [default to undefined]
**completed_at** | **string** |  | [optional] [default to undefined]
**conclusion** | **string** | Conclusion is the outcome of the job. Possible values: \&quot;success\&quot;, \&quot;failure\&quot;, \&quot;neutral\&quot;, \&quot;cancelled\&quot;, \&quot;skipped\&quot;, \&quot;timed_out\&quot;, \&quot;action_required\&quot; | [optional] [default to undefined]
**created_at** | **string** |  | [optional] [default to undefined]
**enterprise_id** | **string** |  | [optional] [default to undefined]
**id** | **number** | ID is the ID of the job. | [optional] [default to undefined]
**labels** | **Array&lt;string&gt;** |  | [optional] [default to undefined]
**locked_by** | **string** |  | [optional] [default to undefined]
**name** | **string** | Name is the name if the job that was triggered. | [optional] [default to undefined]
**org_id** | **string** |  | [optional] [default to undefined]
**repo_id** | **string** | The entity that received the hook.  Webhooks may be configured on the repo, the org and/or the enterprise. If we only configure a repo to use garm, we\&#39;ll only ever receive a webhook from the repo. But if we configure the parent org of the repo and the parent enterprise of the org to use garm, a webhook will be sent for each entity type, in response to one workflow event. Thus, we will get 3 webhooks with the same run_id and job id. Record all involved entities in the same job if we have them configured in garm. | [optional] [default to undefined]
**repository_name** | **string** | repository in which the job was triggered. | [optional] [default to undefined]
**repository_owner** | **string** |  | [optional] [default to undefined]
**run_id** | **number** | RunID is the ID of the workflow run. A run may have multiple jobs. | [optional] [default to undefined]
**runner_group_id** | **number** |  | [optional] [default to undefined]
**runner_group_name** | **string** |  | [optional] [default to undefined]
**runner_id** | **number** |  | [optional] [default to undefined]
**runner_name** | **string** |  | [optional] [default to undefined]
**scaleset_job_id** | **string** | ScaleSetJobID is the job ID when generated for a scale set. | [optional] [default to undefined]
**started_at** | **string** |  | [optional] [default to undefined]
**status** | **string** | Status is the phase of the lifecycle that the job is currently in. \&quot;queued\&quot;, \&quot;in_progress\&quot; and \&quot;completed\&quot;. | [optional] [default to undefined]
**updated_at** | **string** |  | [optional] [default to undefined]
**workflow_job_id** | **number** |  | [optional] [default to undefined]

## Example

```typescript
import { Job } from './api';

const instance: Job = {
    action,
    completed_at,
    conclusion,
    created_at,
    enterprise_id,
    id,
    labels,
    locked_by,
    name,
    org_id,
    repo_id,
    repository_name,
    repository_owner,
    run_id,
    runner_group_id,
    runner_group_name,
    runner_id,
    runner_name,
    scaleset_job_id,
    started_at,
    status,
    updated_at,
    workflow_job_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
