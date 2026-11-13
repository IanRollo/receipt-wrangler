# openapi.model.Receipt

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **String** | Receipt total amount | 
**currency** | **String** | ISO 4217 currency code of the receipt's amount and of its items' amounts (items have no currency of their own). A plain string, never an enum. Never empty on a stored receipt. | [optional] 
**baseAmount** | **String** | What the receipt came to in its group's base currency (decimal string), as the user recorded or the server derived it. Only present when the currency differs from the base currency and a conversion was recorded. | [optional] 
**baseAmountCurrency** | **String** | Read-only, server-set. The ISO 4217 currency baseAmount is in: the group's base currency at the moment the conversion was written. A conversion counts as converted (claim figures, packet) only while this equals the group's current base currency; otherwise the receipt is treated as unconverted and the stored conversion is shown as recorded. Not to be confused with GroupReceiptSettings.baseCurrency. | [optional] 
**exchangeRate** | **String** | Base-currency units per one unit of the receipt's currency (decimal string). Only present alongside baseAmount. | [optional] 
**exchangeRateDate** | **String** | Calendar day (yyyy-MM-dd) the exchange rate applies to | [optional] 
**exchangeRateSource** | **String** | Free text saying where the rate came from, e.g. \"Card statement\" | [optional] 
**conversionNote** | **String** | Read-only warning, present when baseAmount and exchangeRate disagree about the amount by more than 1%. Never blocks a save. | [optional] 
**categories** | [**BuiltList&lt;Category&gt;**](Category.md) | Categories associated to receipt | 
**comments** | [**BuiltList&lt;Comment&gt;**](Comment.md) | Comments associated to receipt | 
**customFields** | [**BuiltList&lt;CustomFieldValue&gt;**](CustomFieldValue.md) | Custom fields associated to receipt | 
**createdAt** | **String** |  | [optional] 
**createdBy** | **int** |  | [optional] 
**date** | **String** | Receipt date | 
**firstComment** | **String** | Text of the receipt's earliest comment the caller may see. Only populated by the paged receipt list (POST /receipt/group/{groupId}); absent when there is no such comment. | [optional] 
**groupId** | **int** | Group foreign key | 
**id** | **int** |  | 
**imageFiles** | [**BuiltList&lt;FileData&gt;**](FileData.md) | Files associated to receipt | [optional] 
**name** | **String** | Receipt name | 
**paidByUserId** | **int** | User paid foreign key | 
**receiptItems** | [**BuiltList&lt;Item&gt;**](Item.md) | Items associated to receipt | 
**resolvedDate** | **String** | Date resolved | [optional] 
**status** | [**ReceiptStatus**](ReceiptStatus.md) |  | 
**tags** | [**BuiltList&lt;Tag&gt;**](Tag.md) | Tags associated to receipt | 
**updatedAt** | **String** |  | [optional] 
**createdByString** | **String** | Created by string, which is anything that is not a user | [optional] 
**emailInbox** | **String** | Set by the server on a receipt created from an email: the System Email login (IMAP username) of the inbox the email was read from. Ignored on write. | [optional] 
**emailSender** | **String** | Set by the server on a receipt created from an email: the sender's From address. Ignored on write. | [optional] 
**claimId** | **int** | Set by the server: the id of the reimbursement claim this receipt is in. Read-only (ignored on write). Omitted when the receipt is in no claim or its claim is not visible to the caller. | [optional] 
**claimName** | **String** | Set by the server: the name of the claim in claimId. Read-only (ignored on write). Omitted when the receipt is in no claim or its claim is not visible to the caller. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


