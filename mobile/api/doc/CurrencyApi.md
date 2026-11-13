# openapi.api.CurrencyApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCurrencies**](CurrencyApi.md#getcurrencies) | **GET** /currency | Get the currency table


# **getCurrencies**
> BuiltList<CurrencyInfo> getCurrencies()

Get the currency table

Returns the active ISO 4217 currency codes the server accepts, with the number of decimal places each allows, ordered by code. [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getCurrencyApi();

try {
    final response = api.getCurrencies();
    print(response);
} catch on DioException (e) {
    print('Exception when calling CurrencyApi->getCurrencies: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BuiltList&lt;CurrencyInfo&gt;**](CurrencyInfo.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

