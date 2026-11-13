# openapi.model.ReceiptSummaryRow

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**status** | [**ReceiptStatus**](ReceiptStatus.md) |  | 
**receiptCount** | **int** | How many receipts this row covers. On the overall row this matches the table's own total count for the same filter. | 
**total** | **String** | Sum of the receipts' amounts | 
**customFieldTotals** | [**BuiltList&lt;ReceiptSummaryCustomFieldTotal&gt;**](ReceiptSummaryCustomFieldTotal.md) | One entry per configured currency custom field, in the configured order. Always present; empty when the group totals only the receipt amount. | 
**currencyTotals** | [**BuiltList&lt;ReceiptSummaryCurrencyTotal&gt;**](ReceiptSummaryCurrencyTotal.md) | The same figures once per currency present in the row (base currency first, then by code), because amounts in different currencies are never added. total and customFieldTotals above cover only receipts in the summary's baseCurrency; receiptCount still counts every receipt. Always present; empty when the row covers no receipt. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


