# openapi.api.ProcessingQueueApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**downloadProcessingQueueTaskSourceFile**](ProcessingQueueApi.md#downloadprocessingqueuetasksourcefile) | **GET** /processingQueue/task/{queue}/{taskId}/sourceFile/download | Download a queued upload&#39;s source file
[**getProcessingQueueForGroup**](ProcessingQueueApi.md#getprocessingqueueforgroup) | **GET** /processingQueue/group/{groupId} | Get the processing queue for a group
[**getProcessingQueueTaskSourceFile**](ProcessingQueueApi.md#getprocessingqueuetasksourcefile) | **GET** /processingQueue/task/{queue}/{taskId}/sourceFile | Get a queued upload&#39;s source file
[**retryProcessingQueueTask**](ProcessingQueueApi.md#retryprocessingqueuetask) | **POST** /processingQueue/task/{queue}/{taskId}/retry | Retry a queued upload now


# **downloadProcessingQueueTaskSourceFile**
> Uint8List downloadProcessingQueueTaskSourceFile(queue, taskId)

Download a queued upload's source file

Returns the upload behind a processing queue item verbatim. Same addressing and gate as getProcessingQueueTaskSourceFile

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getProcessingQueueApi();
final String queue = queue_example; // String | The queue name from a processing queue item, quick_scan or email_receipt_processing
final String taskId = taskId_example; // String | The asynqTaskId of a processing queue item

try {
    final response = api.downloadProcessingQueueTaskSourceFile(queue, taskId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ProcessingQueueApi->downloadProcessingQueueTaskSourceFile: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **queue** | **String**| The queue name from a processing queue item, quick_scan or email_receipt_processing | 
 **taskId** | **String**| The asynqTaskId of a processing queue item | 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProcessingQueueForGroup**
> ProcessingQueueList getProcessingQueueForGroup(groupId)

Get the processing queue for a group

Lists the quick scan and email upload tasks that have not succeeded, read live from the task queue: waiting, being processed, retrying after an error, or failed. Newest first. Each item is judged in its own group, never the group being browsed. A real group is gated on group.activities.read; the synthetic \"All\" group instead lists the items of every group where the caller holds group.activities.read. In a group with member isolation, uploads run by a member the caller cannot see are left out, as in the activity feed. At most 100 tasks are read per queue state; `truncated` says when that cap was hit [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getProcessingQueueApi();
final int groupId = 56; // int | Id of the group, or of the caller's \"All\" group

try {
    final response = api.getProcessingQueueForGroup(groupId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ProcessingQueueApi->getProcessingQueueForGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of the group, or of the caller's \"All\" group | 

### Return type

[**ProcessingQueueList**](ProcessingQueueList.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProcessingQueueTaskSourceFile**
> SystemTaskSourceFileView getProcessingQueueTaskSourceFile(queue, taskId)

Get a queued upload's source file

Returns the upload behind a processing queue item, converted for display. Addressed by queue and asynq task id because a task that has not finished a processing step has no system task for /systemTask/{id}/sourceFile. Gated on group.activities.read in the task's group; an unknown task answers like one in a group the caller cannot reach (403)

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getProcessingQueueApi();
final String queue = queue_example; // String | The queue name from a processing queue item, quick_scan or email_receipt_processing
final String taskId = taskId_example; // String | The asynqTaskId of a processing queue item

try {
    final response = api.getProcessingQueueTaskSourceFile(queue, taskId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ProcessingQueueApi->getProcessingQueueTaskSourceFile: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **queue** | **String**| The queue name from a processing queue item, quick_scan or email_receipt_processing | 
 **taskId** | **String**| The asynqTaskId of a processing queue item | 

### Return type

[**SystemTaskSourceFileView**](SystemTaskSourceFileView.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **retryProcessingQueueTask**
> retryProcessingQueueTask(queue, taskId)

Retry a queued upload now

Runs a RETRYING upload straight away instead of waiting for its next attempt, or reruns a FAILED one whose file is still available (the same rule as rerunActivity). Any other state is a 400. Gated on group.activities.rerun in the task's group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getProcessingQueueApi();
final String queue = queue_example; // String | The queue name from a processing queue item, quick_scan or email_receipt_processing
final String taskId = taskId_example; // String | The asynqTaskId of a processing queue item

try {
    api.retryProcessingQueueTask(queue, taskId);
} catch on DioException (e) {
    print('Exception when calling ProcessingQueueApi->retryProcessingQueueTask: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **queue** | **String**| The queue name from a processing queue item, quick_scan or email_receipt_processing | 
 **taskId** | **String**| The asynqTaskId of a processing queue item | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

