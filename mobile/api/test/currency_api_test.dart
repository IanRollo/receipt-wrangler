import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for CurrencyApi
void main() {
  final instance = Openapi().getCurrencyApi();

  group(CurrencyApi, () {
    // Get the currency table
    //
    // Returns the active ISO 4217 currency codes the server accepts, with the number of decimal places each allows, ordered by code. [SYSTEM USER]
    //
    //Future<BuiltList<CurrencyInfo>> getCurrencies() async
    test('test getCurrencies', () async {
      // TODO
    });

  });
}
