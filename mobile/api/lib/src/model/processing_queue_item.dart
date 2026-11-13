//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'processing_queue_item.g.dart';

/// One quick scan or email upload task, read from the task queue
///
/// Properties:
/// * [asynqTaskId] - The task id; with `queue`, addresses the /processingQueue/task routes
/// * [queue] - The task queue name, e.g. quick_scan or email_receipt_processing. Plain string
/// * [state] - Plain string, not an enum; treat an unknown value as in progress. QUEUED is waiting for a worker, PROCESSING is running now, RETRYING failed an attempt and has another scheduled (see nextRetryAt), FAILED has used its retries and waits for a rerun
/// * [taskType] - QUICK_SCAN or EMAIL_UPLOAD. Plain string, not an enum
/// * [groupId] - The group the upload belongs to
/// * [uploaderUserId] - The user who ran the quick scan. Absent for email uploads
/// * [fileName] - The upload's original file name. Empty for an email with no attachment
/// * [emailSubject] - Email uploads only
/// * [emailSender] - Email uploads only
/// * [uploadedAt] - Modification time of the stored upload, as the task queue does not record when a task was enqueued. Absent once the file is gone
/// * [retryCount] - Failed attempts so far
/// * [maxRetries] - The most retries the queue will schedule after the first attempt
/// * [nextRetryAt] - When the next attempt runs. RETRYING only
/// * [lastError] - The last attempt's error. Empty if none
/// * [lastFailedAt] 
/// * [hasSourceFile] - Whether the upload is still on disk and may be previewed or downloaded
/// * [canBeRestarted] - A FAILED upload whose file is still available, and the caller may rerun it
/// * [canRetryNow] - A RETRYING upload the caller may run now
/// * [systemTaskId] - The system task that /systemTask/{id}/sourceFile and /systemTask/rerunActivity/{id} address. Absent until a processing step has finished and recorded one; use the /processingQueue/task routes, which work for every item
@BuiltValue()
abstract class ProcessingQueueItem implements Built<ProcessingQueueItem, ProcessingQueueItemBuilder> {
  /// The task id; with `queue`, addresses the /processingQueue/task routes
  @BuiltValueField(wireName: r'asynqTaskId')
  String get asynqTaskId;

  /// The task queue name, e.g. quick_scan or email_receipt_processing. Plain string
  @BuiltValueField(wireName: r'queue')
  String get queue;

  /// Plain string, not an enum; treat an unknown value as in progress. QUEUED is waiting for a worker, PROCESSING is running now, RETRYING failed an attempt and has another scheduled (see nextRetryAt), FAILED has used its retries and waits for a rerun
  @BuiltValueField(wireName: r'state')
  String get state;

  /// QUICK_SCAN or EMAIL_UPLOAD. Plain string, not an enum
  @BuiltValueField(wireName: r'taskType')
  String get taskType;

  /// The group the upload belongs to
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  /// The user who ran the quick scan. Absent for email uploads
  @BuiltValueField(wireName: r'uploaderUserId')
  int? get uploaderUserId;

  /// The upload's original file name. Empty for an email with no attachment
  @BuiltValueField(wireName: r'fileName')
  String get fileName;

  /// Email uploads only
  @BuiltValueField(wireName: r'emailSubject')
  String? get emailSubject;

  /// Email uploads only
  @BuiltValueField(wireName: r'emailSender')
  String? get emailSender;

  /// Modification time of the stored upload, as the task queue does not record when a task was enqueued. Absent once the file is gone
  @BuiltValueField(wireName: r'uploadedAt')
  DateTime? get uploadedAt;

  /// Failed attempts so far
  @BuiltValueField(wireName: r'retryCount')
  int get retryCount;

  /// The most retries the queue will schedule after the first attempt
  @BuiltValueField(wireName: r'maxRetries')
  int get maxRetries;

  /// When the next attempt runs. RETRYING only
  @BuiltValueField(wireName: r'nextRetryAt')
  DateTime? get nextRetryAt;

  /// The last attempt's error. Empty if none
  @BuiltValueField(wireName: r'lastError')
  String get lastError;

  @BuiltValueField(wireName: r'lastFailedAt')
  DateTime? get lastFailedAt;

  /// Whether the upload is still on disk and may be previewed or downloaded
  @BuiltValueField(wireName: r'hasSourceFile')
  bool get hasSourceFile;

  /// A FAILED upload whose file is still available, and the caller may rerun it
  @BuiltValueField(wireName: r'canBeRestarted')
  bool get canBeRestarted;

  /// A RETRYING upload the caller may run now
  @BuiltValueField(wireName: r'canRetryNow')
  bool get canRetryNow;

  /// The system task that /systemTask/{id}/sourceFile and /systemTask/rerunActivity/{id} address. Absent until a processing step has finished and recorded one; use the /processingQueue/task routes, which work for every item
  @BuiltValueField(wireName: r'systemTaskId')
  int? get systemTaskId;

  ProcessingQueueItem._();

  factory ProcessingQueueItem([void updates(ProcessingQueueItemBuilder b)]) = _$ProcessingQueueItem;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ProcessingQueueItemBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ProcessingQueueItem> get serializer => _$ProcessingQueueItemSerializer();
}

class _$ProcessingQueueItemSerializer implements PrimitiveSerializer<ProcessingQueueItem> {
  @override
  final Iterable<Type> types = const [ProcessingQueueItem, _$ProcessingQueueItem];

  @override
  final String wireName = r'ProcessingQueueItem';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ProcessingQueueItem object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'asynqTaskId';
    yield serializers.serialize(
      object.asynqTaskId,
      specifiedType: const FullType(String),
    );
    yield r'queue';
    yield serializers.serialize(
      object.queue,
      specifiedType: const FullType(String),
    );
    yield r'state';
    yield serializers.serialize(
      object.state,
      specifiedType: const FullType(String),
    );
    yield r'taskType';
    yield serializers.serialize(
      object.taskType,
      specifiedType: const FullType(String),
    );
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
    if (object.uploaderUserId != null) {
      yield r'uploaderUserId';
      yield serializers.serialize(
        object.uploaderUserId,
        specifiedType: const FullType(int),
      );
    }
    yield r'fileName';
    yield serializers.serialize(
      object.fileName,
      specifiedType: const FullType(String),
    );
    if (object.emailSubject != null) {
      yield r'emailSubject';
      yield serializers.serialize(
        object.emailSubject,
        specifiedType: const FullType(String),
      );
    }
    if (object.emailSender != null) {
      yield r'emailSender';
      yield serializers.serialize(
        object.emailSender,
        specifiedType: const FullType(String),
      );
    }
    if (object.uploadedAt != null) {
      yield r'uploadedAt';
      yield serializers.serialize(
        object.uploadedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    yield r'retryCount';
    yield serializers.serialize(
      object.retryCount,
      specifiedType: const FullType(int),
    );
    yield r'maxRetries';
    yield serializers.serialize(
      object.maxRetries,
      specifiedType: const FullType(int),
    );
    if (object.nextRetryAt != null) {
      yield r'nextRetryAt';
      yield serializers.serialize(
        object.nextRetryAt,
        specifiedType: const FullType(DateTime),
      );
    }
    yield r'lastError';
    yield serializers.serialize(
      object.lastError,
      specifiedType: const FullType(String),
    );
    if (object.lastFailedAt != null) {
      yield r'lastFailedAt';
      yield serializers.serialize(
        object.lastFailedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    yield r'hasSourceFile';
    yield serializers.serialize(
      object.hasSourceFile,
      specifiedType: const FullType(bool),
    );
    yield r'canBeRestarted';
    yield serializers.serialize(
      object.canBeRestarted,
      specifiedType: const FullType(bool),
    );
    yield r'canRetryNow';
    yield serializers.serialize(
      object.canRetryNow,
      specifiedType: const FullType(bool),
    );
    if (object.systemTaskId != null) {
      yield r'systemTaskId';
      yield serializers.serialize(
        object.systemTaskId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    ProcessingQueueItem object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ProcessingQueueItemBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'asynqTaskId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.asynqTaskId = valueDes;
          break;
        case r'queue':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.queue = valueDes;
          break;
        case r'state':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.state = valueDes;
          break;
        case r'taskType':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.taskType = valueDes;
          break;
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'uploaderUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.uploaderUserId = valueDes;
          break;
        case r'fileName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.fileName = valueDes;
          break;
        case r'emailSubject':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.emailSubject = valueDes;
          break;
        case r'emailSender':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.emailSender = valueDes;
          break;
        case r'uploadedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.uploadedAt = valueDes;
          break;
        case r'retryCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.retryCount = valueDes;
          break;
        case r'maxRetries':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.maxRetries = valueDes;
          break;
        case r'nextRetryAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.nextRetryAt = valueDes;
          break;
        case r'lastError':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.lastError = valueDes;
          break;
        case r'lastFailedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.lastFailedAt = valueDes;
          break;
        case r'hasSourceFile':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.hasSourceFile = valueDes;
          break;
        case r'canBeRestarted':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.canBeRestarted = valueDes;
          break;
        case r'canRetryNow':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.canRetryNow = valueDes;
          break;
        case r'systemTaskId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.systemTaskId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ProcessingQueueItem deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ProcessingQueueItemBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

