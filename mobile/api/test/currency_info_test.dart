import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for CurrencyInfo
void main() {
  final instance = CurrencyInfoBuilder();
  // TODO add properties to the builder and call build()

  group(CurrencyInfo, () {
    // ISO 4217 code, e.g. USD. A plain string, never an enum.
    // String code
    test('to test the property `code`', () async {
      // TODO
    });

    // Decimal places an amount in this currency may have (JPY 0, USD 2, JOD 3). A more precise amount is refused, never rounded.
    // int exponent
    test('to test the property `exponent`', () async {
      // TODO
    });

    // Display symbol where unambiguous; absent otherwise
    // String symbol
    test('to test the property `symbol`', () async {
      // TODO
    });

  });
}
