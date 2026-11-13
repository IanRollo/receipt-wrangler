// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'processing_queue_item.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ProcessingQueueItem extends ProcessingQueueItem {
  @override
  final String asynqTaskId;
  @override
  final String queue;
  @override
  final String state;
  @override
  final String taskType;
  @override
  final int groupId;
  @override
  final int? uploaderUserId;
  @override
  final String fileName;
  @override
  final String? emailSubject;
  @override
  final String? emailSender;
  @override
  final DateTime? uploadedAt;
  @override
  final int retryCount;
  @override
  final int maxRetries;
  @override
  final DateTime? nextRetryAt;
  @override
  final String lastError;
  @override
  final DateTime? lastFailedAt;
  @override
  final bool hasSourceFile;
  @override
  final bool canBeRestarted;
  @override
  final bool canRetryNow;
  @override
  final int? systemTaskId;

  factory _$ProcessingQueueItem(
          [void Function(ProcessingQueueItemBuilder)? updates]) =>
      (ProcessingQueueItemBuilder()..update(updates))._build();

  _$ProcessingQueueItem._(
      {required this.asynqTaskId,
      required this.queue,
      required this.state,
      required this.taskType,
      required this.groupId,
      this.uploaderUserId,
      required this.fileName,
      this.emailSubject,
      this.emailSender,
      this.uploadedAt,
      required this.retryCount,
      required this.maxRetries,
      this.nextRetryAt,
      required this.lastError,
      this.lastFailedAt,
      required this.hasSourceFile,
      required this.canBeRestarted,
      required this.canRetryNow,
      this.systemTaskId})
      : super._();
  @override
  ProcessingQueueItem rebuild(
          void Function(ProcessingQueueItemBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ProcessingQueueItemBuilder toBuilder() =>
      ProcessingQueueItemBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ProcessingQueueItem &&
        asynqTaskId == other.asynqTaskId &&
        queue == other.queue &&
        state == other.state &&
        taskType == other.taskType &&
        groupId == other.groupId &&
        uploaderUserId == other.uploaderUserId &&
        fileName == other.fileName &&
        emailSubject == other.emailSubject &&
        emailSender == other.emailSender &&
        uploadedAt == other.uploadedAt &&
        retryCount == other.retryCount &&
        maxRetries == other.maxRetries &&
        nextRetryAt == other.nextRetryAt &&
        lastError == other.lastError &&
        lastFailedAt == other.lastFailedAt &&
        hasSourceFile == other.hasSourceFile &&
        canBeRestarted == other.canBeRestarted &&
        canRetryNow == other.canRetryNow &&
        systemTaskId == other.systemTaskId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, asynqTaskId.hashCode);
    _$hash = $jc(_$hash, queue.hashCode);
    _$hash = $jc(_$hash, state.hashCode);
    _$hash = $jc(_$hash, taskType.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, uploaderUserId.hashCode);
    _$hash = $jc(_$hash, fileName.hashCode);
    _$hash = $jc(_$hash, emailSubject.hashCode);
    _$hash = $jc(_$hash, emailSender.hashCode);
    _$hash = $jc(_$hash, uploadedAt.hashCode);
    _$hash = $jc(_$hash, retryCount.hashCode);
    _$hash = $jc(_$hash, maxRetries.hashCode);
    _$hash = $jc(_$hash, nextRetryAt.hashCode);
    _$hash = $jc(_$hash, lastError.hashCode);
    _$hash = $jc(_$hash, lastFailedAt.hashCode);
    _$hash = $jc(_$hash, hasSourceFile.hashCode);
    _$hash = $jc(_$hash, canBeRestarted.hashCode);
    _$hash = $jc(_$hash, canRetryNow.hashCode);
    _$hash = $jc(_$hash, systemTaskId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ProcessingQueueItem')
          ..add('asynqTaskId', asynqTaskId)
          ..add('queue', queue)
          ..add('state', state)
          ..add('taskType', taskType)
          ..add('groupId', groupId)
          ..add('uploaderUserId', uploaderUserId)
          ..add('fileName', fileName)
          ..add('emailSubject', emailSubject)
          ..add('emailSender', emailSender)
          ..add('uploadedAt', uploadedAt)
          ..add('retryCount', retryCount)
          ..add('maxRetries', maxRetries)
          ..add('nextRetryAt', nextRetryAt)
          ..add('lastError', lastError)
          ..add('lastFailedAt', lastFailedAt)
          ..add('hasSourceFile', hasSourceFile)
          ..add('canBeRestarted', canBeRestarted)
          ..add('canRetryNow', canRetryNow)
          ..add('systemTaskId', systemTaskId))
        .toString();
  }
}

class ProcessingQueueItemBuilder
    implements Builder<ProcessingQueueItem, ProcessingQueueItemBuilder> {
  _$ProcessingQueueItem? _$v;

  String? _asynqTaskId;
  String? get asynqTaskId => _$this._asynqTaskId;
  set asynqTaskId(String? asynqTaskId) => _$this._asynqTaskId = asynqTaskId;

  String? _queue;
  String? get queue => _$this._queue;
  set queue(String? queue) => _$this._queue = queue;

  String? _state;
  String? get state => _$this._state;
  set state(String? state) => _$this._state = state;

  String? _taskType;
  String? get taskType => _$this._taskType;
  set taskType(String? taskType) => _$this._taskType = taskType;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _uploaderUserId;
  int? get uploaderUserId => _$this._uploaderUserId;
  set uploaderUserId(int? uploaderUserId) =>
      _$this._uploaderUserId = uploaderUserId;

  String? _fileName;
  String? get fileName => _$this._fileName;
  set fileName(String? fileName) => _$this._fileName = fileName;

  String? _emailSubject;
  String? get emailSubject => _$this._emailSubject;
  set emailSubject(String? emailSubject) => _$this._emailSubject = emailSubject;

  String? _emailSender;
  String? get emailSender => _$this._emailSender;
  set emailSender(String? emailSender) => _$this._emailSender = emailSender;

  DateTime? _uploadedAt;
  DateTime? get uploadedAt => _$this._uploadedAt;
  set uploadedAt(DateTime? uploadedAt) => _$this._uploadedAt = uploadedAt;

  int? _retryCount;
  int? get retryCount => _$this._retryCount;
  set retryCount(int? retryCount) => _$this._retryCount = retryCount;

  int? _maxRetries;
  int? get maxRetries => _$this._maxRetries;
  set maxRetries(int? maxRetries) => _$this._maxRetries = maxRetries;

  DateTime? _nextRetryAt;
  DateTime? get nextRetryAt => _$this._nextRetryAt;
  set nextRetryAt(DateTime? nextRetryAt) => _$this._nextRetryAt = nextRetryAt;

  String? _lastError;
  String? get lastError => _$this._lastError;
  set lastError(String? lastError) => _$this._lastError = lastError;

  DateTime? _lastFailedAt;
  DateTime? get lastFailedAt => _$this._lastFailedAt;
  set lastFailedAt(DateTime? lastFailedAt) =>
      _$this._lastFailedAt = lastFailedAt;

  bool? _hasSourceFile;
  bool? get hasSourceFile => _$this._hasSourceFile;
  set hasSourceFile(bool? hasSourceFile) =>
      _$this._hasSourceFile = hasSourceFile;

  bool? _canBeRestarted;
  bool? get canBeRestarted => _$this._canBeRestarted;
  set canBeRestarted(bool? canBeRestarted) =>
      _$this._canBeRestarted = canBeRestarted;

  bool? _canRetryNow;
  bool? get canRetryNow => _$this._canRetryNow;
  set canRetryNow(bool? canRetryNow) => _$this._canRetryNow = canRetryNow;

  int? _systemTaskId;
  int? get systemTaskId => _$this._systemTaskId;
  set systemTaskId(int? systemTaskId) => _$this._systemTaskId = systemTaskId;

  ProcessingQueueItemBuilder() {
    ProcessingQueueItem._defaults(this);
  }

  ProcessingQueueItemBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _asynqTaskId = $v.asynqTaskId;
      _queue = $v.queue;
      _state = $v.state;
      _taskType = $v.taskType;
      _groupId = $v.groupId;
      _uploaderUserId = $v.uploaderUserId;
      _fileName = $v.fileName;
      _emailSubject = $v.emailSubject;
      _emailSender = $v.emailSender;
      _uploadedAt = $v.uploadedAt;
      _retryCount = $v.retryCount;
      _maxRetries = $v.maxRetries;
      _nextRetryAt = $v.nextRetryAt;
      _lastError = $v.lastError;
      _lastFailedAt = $v.lastFailedAt;
      _hasSourceFile = $v.hasSourceFile;
      _canBeRestarted = $v.canBeRestarted;
      _canRetryNow = $v.canRetryNow;
      _systemTaskId = $v.systemTaskId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ProcessingQueueItem other) {
    _$v = other as _$ProcessingQueueItem;
  }

  @override
  void update(void Function(ProcessingQueueItemBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ProcessingQueueItem build() => _build();

  _$ProcessingQueueItem _build() {
    final _$result = _$v ??
        _$ProcessingQueueItem._(
          asynqTaskId: BuiltValueNullFieldError.checkNotNull(
              asynqTaskId, r'ProcessingQueueItem', 'asynqTaskId'),
          queue: BuiltValueNullFieldError.checkNotNull(
              queue, r'ProcessingQueueItem', 'queue'),
          state: BuiltValueNullFieldError.checkNotNull(
              state, r'ProcessingQueueItem', 'state'),
          taskType: BuiltValueNullFieldError.checkNotNull(
              taskType, r'ProcessingQueueItem', 'taskType'),
          groupId: BuiltValueNullFieldError.checkNotNull(
              groupId, r'ProcessingQueueItem', 'groupId'),
          uploaderUserId: uploaderUserId,
          fileName: BuiltValueNullFieldError.checkNotNull(
              fileName, r'ProcessingQueueItem', 'fileName'),
          emailSubject: emailSubject,
          emailSender: emailSender,
          uploadedAt: uploadedAt,
          retryCount: BuiltValueNullFieldError.checkNotNull(
              retryCount, r'ProcessingQueueItem', 'retryCount'),
          maxRetries: BuiltValueNullFieldError.checkNotNull(
              maxRetries, r'ProcessingQueueItem', 'maxRetries'),
          nextRetryAt: nextRetryAt,
          lastError: BuiltValueNullFieldError.checkNotNull(
              lastError, r'ProcessingQueueItem', 'lastError'),
          lastFailedAt: lastFailedAt,
          hasSourceFile: BuiltValueNullFieldError.checkNotNull(
              hasSourceFile, r'ProcessingQueueItem', 'hasSourceFile'),
          canBeRestarted: BuiltValueNullFieldError.checkNotNull(
              canBeRestarted, r'ProcessingQueueItem', 'canBeRestarted'),
          canRetryNow: BuiltValueNullFieldError.checkNotNull(
              canRetryNow, r'ProcessingQueueItem', 'canRetryNow'),
          systemTaskId: systemTaskId,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
