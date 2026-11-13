# openapi.model.ProcessingQueueItem

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**asynqTaskId** | **String** | The task id; with `queue`, addresses the /processingQueue/task routes | 
**queue** | **String** | The task queue name, e.g. quick_scan or email_receipt_processing. Plain string | 
**state** | **String** | Plain string, not an enum; treat an unknown value as in progress. QUEUED is waiting for a worker, PROCESSING is running now, RETRYING failed an attempt and has another scheduled (see nextRetryAt), FAILED has used its retries and waits for a rerun | 
**taskType** | **String** | QUICK_SCAN or EMAIL_UPLOAD. Plain string, not an enum | 
**groupId** | **int** | The group the upload belongs to | 
**uploaderUserId** | **int** | The user who ran the quick scan. Absent for email uploads | [optional] 
**fileName** | **String** | The upload's original file name. Empty for an email with no attachment | 
**emailSubject** | **String** | Email uploads only | [optional] 
**emailSender** | **String** | Email uploads only | [optional] 
**uploadedAt** | [**DateTime**](DateTime.md) | Modification time of the stored upload, as the task queue does not record when a task was enqueued. Absent once the file is gone | [optional] 
**retryCount** | **int** | Failed attempts so far | 
**maxRetries** | **int** | The most retries the queue will schedule after the first attempt | 
**nextRetryAt** | [**DateTime**](DateTime.md) | When the next attempt runs. RETRYING only | [optional] 
**lastError** | **String** | The last attempt's error. Empty if none | 
**lastFailedAt** | [**DateTime**](DateTime.md) |  | [optional] 
**hasSourceFile** | **bool** | Whether the upload is still on disk and may be previewed or downloaded | 
**canBeRestarted** | **bool** | A FAILED upload whose file is still available, and the caller may rerun it | 
**canRetryNow** | **bool** | A RETRYING upload the caller may run now | 
**systemTaskId** | **int** | The system task that /systemTask/{id}/sourceFile and /systemTask/rerunActivity/{id} address. Absent until a processing step has finished and recorded one; use the /processingQueue/task routes, which work for every item | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


