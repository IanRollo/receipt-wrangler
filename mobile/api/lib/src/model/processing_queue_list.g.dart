// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'processing_queue_list.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ProcessingQueueList extends ProcessingQueueList {
  @override
  final BuiltList<ProcessingQueueItem> items;
  @override
  final bool truncated;

  factory _$ProcessingQueueList(
          [void Function(ProcessingQueueListBuilder)? updates]) =>
      (ProcessingQueueListBuilder()..update(updates))._build();

  _$ProcessingQueueList._({required this.items, required this.truncated})
      : super._();
  @override
  ProcessingQueueList rebuild(
          void Function(ProcessingQueueListBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ProcessingQueueListBuilder toBuilder() =>
      ProcessingQueueListBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ProcessingQueueList &&
        items == other.items &&
        truncated == other.truncated;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, items.hashCode);
    _$hash = $jc(_$hash, truncated.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ProcessingQueueList')
          ..add('items', items)
          ..add('truncated', truncated))
        .toString();
  }
}

class ProcessingQueueListBuilder
    implements Builder<ProcessingQueueList, ProcessingQueueListBuilder> {
  _$ProcessingQueueList? _$v;

  ListBuilder<ProcessingQueueItem>? _items;
  ListBuilder<ProcessingQueueItem> get items =>
      _$this._items ??= ListBuilder<ProcessingQueueItem>();
  set items(ListBuilder<ProcessingQueueItem>? items) => _$this._items = items;

  bool? _truncated;
  bool? get truncated => _$this._truncated;
  set truncated(bool? truncated) => _$this._truncated = truncated;

  ProcessingQueueListBuilder() {
    ProcessingQueueList._defaults(this);
  }

  ProcessingQueueListBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _items = $v.items.toBuilder();
      _truncated = $v.truncated;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ProcessingQueueList other) {
    _$v = other as _$ProcessingQueueList;
  }

  @override
  void update(void Function(ProcessingQueueListBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ProcessingQueueList build() => _build();

  _$ProcessingQueueList _build() {
    _$ProcessingQueueList _$result;
    try {
      _$result = _$v ??
          _$ProcessingQueueList._(
            items: items.build(),
            truncated: BuiltValueNullFieldError.checkNotNull(
                truncated, r'ProcessingQueueList', 'truncated'),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'items';
        items.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ProcessingQueueList', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
