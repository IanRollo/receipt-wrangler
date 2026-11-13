//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/receipt_summary_custom_field_total.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt_summary_currency_total.g.dart';

/// A summary row's figures for one currency
///
/// Properties:
/// * [currency] - ISO 4217 code
/// * [receiptCount] - Receipts in this currency
/// * [total] - Sum of those receipts' amounts, in this currency
/// * [customFieldTotals] - The configured currency custom fields' totals over those receipts, in the same order as the row's customFieldTotals. Values are read in the receipt's currency.
@BuiltValue()
abstract class ReceiptSummaryCurrencyTotal implements Built<ReceiptSummaryCurrencyTotal, ReceiptSummaryCurrencyTotalBuilder> {
  /// ISO 4217 code
  @BuiltValueField(wireName: r'currency')
  String get currency;

  /// Receipts in this currency
  @BuiltValueField(wireName: r'receiptCount')
  int get receiptCount;

  /// Sum of those receipts' amounts, in this currency
  @BuiltValueField(wireName: r'total')
  String get total;

  /// The configured currency custom fields' totals over those receipts, in the same order as the row's customFieldTotals. Values are read in the receipt's currency.
  @BuiltValueField(wireName: r'customFieldTotals')
  BuiltList<ReceiptSummaryCustomFieldTotal> get customFieldTotals;

  ReceiptSummaryCurrencyTotal._();

  factory ReceiptSummaryCurrencyTotal([void updates(ReceiptSummaryCurrencyTotalBuilder b)]) = _$ReceiptSummaryCurrencyTotal;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptSummaryCurrencyTotalBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<ReceiptSummaryCurrencyTotal> get serializer => _$ReceiptSummaryCurrencyTotalSerializer();
}

class _$ReceiptSummaryCurrencyTotalSerializer implements PrimitiveSerializer<ReceiptSummaryCurrencyTotal> {
  @override
  final Iterable<Type> types = const [ReceiptSummaryCurrencyTotal, _$ReceiptSummaryCurrencyTotal];

  @override
  final String wireName = r'ReceiptSummaryCurrencyTotal';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    ReceiptSummaryCurrencyTotal object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'currency';
    yield serializers.serialize(
      object.currency,
      specifiedType: const FullType(String),
    );
    yield r'receiptCount';
    yield serializers.serialize(
      object.receiptCount,
      specifiedType: const FullType(int),
    );
    yield r'total';
    yield serializers.serialize(
      object.total,
      specifiedType: const FullType(String),
    );
    yield r'customFieldTotals';
    yield serializers.serialize(
      object.customFieldTotals,
      specifiedType: const FullType(BuiltList, [FullType(ReceiptSummaryCustomFieldTotal)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    ReceiptSummaryCurrencyTotal object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptSummaryCurrencyTotalBuilder result,
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
        case r'receiptCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptCount = valueDes;
          break;
        case r'total':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.total = valueDes;
          break;
        case r'customFieldTotals':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(ReceiptSummaryCustomFieldTotal)]),
          ) as BuiltList<ReceiptSummaryCustomFieldTotal>;
          result.customFieldTotals.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  ReceiptSummaryCurrencyTotal deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptSummaryCurrencyTotalBuilder();
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

