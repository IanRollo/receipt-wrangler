# openapi.model.ReceiptSummaryCurrencyTotal

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**currency** | **String** | ISO 4217 code | 
**receiptCount** | **int** | Receipts in this currency | 
**total** | **String** | Sum of those receipts' amounts, in this currency | 
**customFieldTotals** | [**BuiltList&lt;ReceiptSummaryCustomFieldTotal&gt;**](ReceiptSummaryCustomFieldTotal.md) | The configured currency custom fields' totals over those receipts, in the same order as the row's customFieldTotals. Values are read in the receipt's currency. | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


