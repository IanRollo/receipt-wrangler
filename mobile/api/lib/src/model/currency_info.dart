//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'currency_info.g.dart';

/// One currency in the server's ISO 4217 table
///
/// Properties:
/// * [code] - ISO 4217 code, e.g. USD. A plain string, never an enum.
/// * [exponent] - Decimal places an amount in this currency may have (JPY 0, USD 2, JOD 3). A more precise amount is refused, never rounded.
/// * [symbol] - Display symbol where unambiguous; absent otherwise
@BuiltValue()
abstract class CurrencyInfo implements Built<CurrencyInfo, CurrencyInfoBuilder> {
  /// ISO 4217 code, e.g. USD. A plain string, never an enum.
  @BuiltValueField(wireName: r'code')
  String get code;

  /// Decimal places an amount in this currency may have (JPY 0, USD 2, JOD 3). A more precise amount is refused, never rounded.
  @BuiltValueField(wireName: r'exponent')
  int get exponent;

  /// Display symbol where unambiguous; absent otherwise
  @BuiltValueField(wireName: r'symbol')
  String? get symbol;

  CurrencyInfo._();

  factory CurrencyInfo([void updates(CurrencyInfoBuilder b)]) = _$CurrencyInfo;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CurrencyInfoBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<CurrencyInfo> get serializer => _$CurrencyInfoSerializer();
}

class _$CurrencyInfoSerializer implements PrimitiveSerializer<CurrencyInfo> {
  @override
  final Iterable<Type> types = const [CurrencyInfo, _$CurrencyInfo];

  @override
  final String wireName = r'CurrencyInfo';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CurrencyInfo object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'code';
    yield serializers.serialize(
      object.code,
      specifiedType: const FullType(String),
    );
    yield r'exponent';
    yield serializers.serialize(
      object.exponent,
      specifiedType: const FullType(int),
    );
    if (object.symbol != null) {
      yield r'symbol';
      yield serializers.serialize(
        object.symbol,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    CurrencyInfo object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CurrencyInfoBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'code':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.code = valueDes;
          break;
        case r'exponent':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.exponent = valueDes;
          break;
        case r'symbol':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.symbol = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  CurrencyInfo deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CurrencyInfoBuilder();
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

