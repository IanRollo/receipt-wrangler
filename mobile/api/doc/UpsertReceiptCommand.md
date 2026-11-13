# openapi.model.UpsertReceiptCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Receipt name | 
**amount** | **String** | Receipt total amount | 
**currency** | **String** | ISO 4217 currency code. Empty or absent means the group's default currency on a create and the receipt's stored currency on an update. An unknown code is a 400 under the currency key. Amounts with more decimal places than the currency allows are a 400 under the amount key (JPY whole, USD 2, JOD 3). | [optional] 
**baseAmount** | **String** | What the receipt came to in the group's base currency. Optional; the server derives it from exchangeRate when only that is given. Ignored (cleared) when currency equals the base currency. On an update, omitting all four conversion keys keeps the stored conversion; sending any of them replaces the whole conversion. | [optional] 
**exchangeRate** | **String** | Base-currency units per one unit of currency. Optional; derived from baseAmount (6 decimal places) when only that is given. Send null or omit it, never \"\". | [optional] 
**exchangeRateDate** | **String** | Calendar day (yyyy-MM-dd) the rate applies to. Optional. | [optional] 
**exchangeRateSource** | **String** | Free text, at most 255 characters. Optional. | [optional] 
**date** | **String** | Receipt date | 
**groupId** | **int** | Group foreign key | 
**paidByUserId** | **int** | User paid foreign key | 
**status** | [**ReceiptStatus**](ReceiptStatus.md) |  | 
**categories** | [**BuiltList&lt;UpsertCategoryCommand&gt;**](UpsertCategoryCommand.md) | Categories associated to receipt | [optional] 
**tags** | [**BuiltList&lt;UpsertTagCommand&gt;**](UpsertTagCommand.md) | Tags associated to receipt | [optional] 
**receiptItems** | [**BuiltList&lt;UpsertItemCommand&gt;**](UpsertItemCommand.md) | Items associated to receipt | [optional] 
**comments** | [**BuiltList&lt;UpsertCommentCommand&gt;**](UpsertCommentCommand.md) | Comments associated to receipt | [optional] 
**customFields** | [**BuiltList&lt;UpsertCustomFieldValueCommand&gt;**](UpsertCustomFieldValueCommand.md) | Custom fields associated to receipt | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


