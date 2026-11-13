import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptSummaryCurrencyTotal
void main() {
  final instance = ReceiptSummaryCurrencyTotalBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptSummaryCurrencyTotal, () {
    // ISO 4217 code
    // String currency
    test('to test the property `currency`', () async {
      // TODO
    });

    // Receipts in this currency
    // int receiptCount
    test('to test the property `receiptCount`', () async {
      // TODO
    });

    // Sum of those receipts' amounts, in this currency
    // String total
    test('to test the property `total`', () async {
      // TODO
    });

    // The configured currency custom fields' totals over those receipts, in the same order as the row's customFieldTotals. Values are read in the receipt's currency.
    // BuiltList<ReceiptSummaryCustomFieldTotal> customFieldTotals
    test('to test the property `customFieldTotals`', () async {
      // TODO
    });

  });
}
