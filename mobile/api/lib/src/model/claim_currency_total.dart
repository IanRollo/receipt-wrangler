//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'claim_currency_total.g.dart';

/// An amount in one currency
///
/// Properties:
/// * [currency] - ISO 4217 code
/// * [total] - Decimal string at the currency's own number of decimal places
@BuiltValue()
abstract class ClaimCurrencyTotal implements Built<ClaimCurrencyTotal, ClaimCurrencyTotalBuilder> {
  /// ISO 4217 code
  @BuiltValueField(wireName: r'currency')
  String get currency;

  /// Decimal string at the currency's own number of decimal places
  @BuiltValueField(wireName: r'total')
  String get total;

  ClaimCurrencyTotal._();

  factory ClaimCurrencyTotal([void updates(ClaimCurrencyTotalBuilder b)]) = _$ClaimCurrencyTotal;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ClaimCurrencyTotalBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ClaimCurrencyTotal> get serializer => _$ClaimCurrencyTotalSerializer();
}

class _$ClaimCurrencyTotalSerializer implements PrimitiveSerializer<ClaimCurrencyTotal> {
  @override
  final Iterable<Type> types = const [ClaimCurrencyTotal, _$ClaimCurrencyTotal];

  @override
  final String wireName = r'ClaimCurrencyTotal';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ClaimCurrencyTotal object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'currency';
    yield serializers.serialize(
      object.currency,
      specifiedType: const FullType(String),
    );
    yield r'total';
    yield serializers.serialize(
      object.total,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ClaimCurrencyTotal object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ClaimCurrencyTotalBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'currency':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.currency = valueDes;
          break;
        case r'total':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.total = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ClaimCurrencyTotal deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ClaimCurrencyTotalBuilder();
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

