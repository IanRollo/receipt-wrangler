import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ProcessingQueueItem
void main() {
  final instance = ProcessingQueueItemBuilder();
  // TODO add properties to the builder and call build()

  group(ProcessingQueueItem, () {
    // The task id; with `queue`, addresses the /processingQueue/task routes
    // String asynqTaskId
    test('to test the property `asynqTaskId`', () async {
      // TODO
    });

    // The task queue name, e.g. quick_scan or email_receipt_processing. Plain string
    // String queue
    test('to test the property `queue`', () async {
      // TODO
    });

    // Plain string, not an enum; treat an unknown value as in progress. QUEUED is waiting for a worker, PROCESSING is running now, RETRYING failed an attempt and has another scheduled (see nextRetryAt), FAILED has used its retries and waits for a rerun
    // String state
    test('to test the property `state`', () async {
      // TODO
    });

    // QUICK_SCAN or EMAIL_UPLOAD. Plain string, not an enum
    // String taskType
    test('to test the property `taskType`', () async {
      // TODO
    });

    // The group the upload belongs to
    // int groupId
    test('to test the property `groupId`', () async {
      // TODO
    });

    // The user who ran the quick scan. Absent for email uploads
    // int uploaderUserId
    test('to test the property `uploaderUserId`', () async {
      // TODO
    });

    // The upload's original file name. Empty for an email with no attachment
    // String fileName
    test('to test the property `fileName`', () async {
      // TODO
    });

    // Email uploads only
    // String emailSubject
    test('to test the property `emailSubject`', () async {
      // TODO
    });

    // Email uploads only
    // String emailSender
    test('to test the property `emailSender`', () async {
      // TODO
    });

    // Modification time of the stored upload, as the task queue does not record when a task was enqueued. Absent once the file is gone
    // DateTime uploadedAt
    test('to test the property `uploadedAt`', () async {
      // TODO
    });

    // Failed attempts so far
    // int retryCount
    test('to test the property `retryCount`', () async {
      // TODO
    });

    // The most retries the queue will schedule after the first attempt
    // int maxRetries
    test('to test the property `maxRetries`', () async {
      // TODO
    });

    // When the next attempt runs. RETRYING only
    // DateTime nextRetryAt
    test('to test the property `nextRetryAt`', () async {
      // TODO
    });

    // The last attempt's error. Empty if none
    // String lastError
    test('to test the property `lastError`', () async {
      // TODO
    });

    // DateTime lastFailedAt
    test('to test the property `lastFailedAt`', () async {
      // TODO
    });

    // Whether the upload is still on disk and may be previewed or downloaded
    // bool hasSourceFile
    test('to test the property `hasSourceFile`', () async {
      // TODO
    });

    // A FAILED upload whose file is still available, and the caller may rerun it
    // bool canBeRestarted
    test('to test the property `canBeRestarted`', () async {
      // TODO
    });

    // A RETRYING upload the caller may run now
    // bool canRetryNow
    test('to test the property `canRetryNow`', () async {
      // TODO
    });

    // The system task that /systemTask/{id}/sourceFile and /systemTask/rerunActivity/{id} address. Absent until a processing step has finished and recorded one; use the /processingQueue/task routes, which work for every item
    // int systemTaskId
    test('to test the property `systemTaskId`', () async {
      // TODO
    });

  });
}
