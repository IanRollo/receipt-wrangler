import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ProcessingQueueApi
void main() {
  final instance = Openapi().getProcessingQueueApi();

  group(ProcessingQueueApi, () {
    // Download a queued upload's source file
    //
    // Returns the upload behind a processing queue item verbatim. Same addressing and gate as getProcessingQueueTaskSourceFile
    //
    //Future<Uint8List> downloadProcessingQueueTaskSourceFile(String queue, String taskId) async
    test('test downloadProcessingQueueTaskSourceFile', () async {
      // TODO
    });

    // Get the processing queue for a group
    //
    // Lists the quick scan and email upload tasks that have not succeeded, read live from the task queue: waiting, being processed, retrying after an error, or failed. Newest first. Each item is judged in its own group, never the group being browsed. A real group is gated on group.activities.read; the synthetic \"All\" group instead lists the items of every group where the caller holds group.activities.read. In a group with member isolation, uploads run by a member the caller cannot see are left out, as in the activity feed. At most 100 tasks are read per queue state; `truncated` says when that cap was hit [SYSTEM USER]
    //
    //Future<ProcessingQueueList> getProcessingQueueForGroup(int groupId) async
    test('test getProcessingQueueForGroup', () async {
      // TODO
    });

    // Get a queued upload's source file
    //
    // Returns the upload behind a processing queue item, converted for display. Addressed by queue and asynq task id because a task that has not finished a processing step has no system task for /systemTask/{id}/sourceFile. Gated on group.activities.read in the task's group; an unknown task answers like one in a group the caller cannot reach (403)
    //
    //Future<SystemTaskSourceFileView> getProcessingQueueTaskSourceFile(String queue, String taskId) async
    test('test getProcessingQueueTaskSourceFile', () async {
      // TODO
    });

    // Retry a queued upload now
    //
    // Runs a RETRYING upload straight away instead of waiting for its next attempt, or reruns a FAILED one whose file is still available (the same rule as rerunActivity). Any other state is a 400. Gated on group.activities.rerun in the task's group
    //
    //Future retryProcessingQueueTask(String queue, String taskId) async
    test('test retryProcessingQueueTask', () async {
      // TODO
    });

  });
}
