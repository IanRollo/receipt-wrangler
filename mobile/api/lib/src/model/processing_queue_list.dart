//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/processing_queue_item.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'processing_queue_list.g.dart';

/// The quick scan and email upload tasks that have not succeeded
///
/// Properties:
/// * [items] - Newest first. Always present; empty when none.
/// * [truncated] - True when some queue state held more tasks than the 100 read, so the list may omit some
@BuiltValue()
abstract class ProcessingQueueList implements Built<ProcessingQueueList, ProcessingQueueListBuilder> {
  /// Newest first. Always present; empty when none.
  @BuiltValueField(wireName: r'items')
  BuiltList<ProcessingQueueItem> get items;

  /// True when some queue state held more tasks than the 100 read, so the list may omit some
  @BuiltValueField(wireName: r'truncated')
  bool get truncated;

  ProcessingQueueList._();

  factory ProcessingQueueList([void updates(ProcessingQueueListBuilder b)]) = _$ProcessingQueueList;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ProcessingQueueListBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ProcessingQueueList> get serializer => _$ProcessingQueueListSerializer();
}

class _$ProcessingQueueListSerializer implements PrimitiveSerializer<ProcessingQueueList> {
  @override
  final Iterable<Type> types = const [ProcessingQueueList, _$ProcessingQueueList];

  @override
  final String wireName = r'ProcessingQueueList';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ProcessingQueueList object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'items';
    yield serializers.serialize(
      object.items,
      specifiedType: const FullType(BuiltList, [FullType(ProcessingQueueItem)]),
    );
    yield r'truncated';
    yield serializers.serialize(
      object.truncated,
      specifiedType: const FullType(bool),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ProcessingQueueList object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ProcessingQueueListBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'items':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ProcessingQueueItem)]),
          ) as BuiltList<ProcessingQueueItem>;
          result.items.replace(valueDes);
          break;
        case r'truncated':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.truncated = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ProcessingQueueList deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ProcessingQueueListBuilder();
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

