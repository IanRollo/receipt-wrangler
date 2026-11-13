import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ProcessingQueueList
void main() {
  final instance = ProcessingQueueListBuilder();
  // TODO add properties to the builder and call build()

  group(ProcessingQueueList, () {
    // Newest first. Always present; empty when none.
    // BuiltList<ProcessingQueueItem> items
    test('to test the property `items`', () async {
      // TODO
    });

    // True when some queue state held more tasks than the 100 read, so the list may omit some
    // bool truncated
    test('to test the property `truncated`', () async {
      // TODO
    });

  });
}
