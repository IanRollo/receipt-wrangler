//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_status.dart';
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/category.dart';
import 'package:openapi/src/model/custom_field_value.dart';
import 'package:openapi/src/model/tag.dart';
import 'package:openapi/src/model/comment.dart';
import 'package:openapi/src/model/file_data.dart';
import 'package:openapi/src/model/item.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'receipt.g.dart';

/// Receipt
///
/// Properties:
/// * [amount] - Receipt total amount
/// * [currency] - ISO 4217 currency code of the receipt's amount and of its items' amounts (items have no currency of their own). A plain string, never an enum. Never empty on a stored receipt.
/// * [baseAmount] - What the receipt came to in its group's base currency (decimal string), as the user recorded or the server derived it. Only present when the currency differs from the base currency and a conversion was recorded.
/// * [baseAmountCurrency] - Read-only, server-set. The ISO 4217 currency baseAmount is in: the group's base currency at the moment the conversion was written. A conversion counts as converted (claim figures, packet) only while this equals the group's current base currency; otherwise the receipt is treated as unconverted and the stored conversion is shown as recorded. Not to be confused with GroupReceiptSettings.baseCurrency.
/// * [exchangeRate] - Base-currency units per one unit of the receipt's currency (decimal string). Only present alongside baseAmount.
/// * [exchangeRateDate] - Calendar day (yyyy-MM-dd) the exchange rate applies to
/// * [exchangeRateSource] - Free text saying where the rate came from, e.g. \"Card statement\"
/// * [conversionNote] - Read-only warning, present when baseAmount and exchangeRate disagree about the amount by more than 1%. Never blocks a save.
/// * [categories] - Categories associated to receipt
/// * [comments] - Comments associated to receipt
/// * [customFields] - Custom fields associated to receipt
/// * [createdAt] 
/// * [createdBy] 
/// * [date] - Receipt date
/// * [firstComment] - Text of the receipt's earliest comment the caller may see. Only populated by the paged receipt list (POST /receipt/group/{groupId}); absent when there is no such comment.
/// * [groupId] - Group foreign key
/// * [id] 
/// * [imageFiles] - Files associated to receipt
/// * [name] - Receipt name
/// * [paidByUserId] - User paid foreign key
/// * [receiptItems] - Items associated to receipt
/// * [resolvedDate] - Date resolved
/// * [status] 
/// * [tags] - Tags associated to receipt
/// * [updatedAt] 
/// * [createdByString] - Created by string, which is anything that is not a user
/// * [emailInbox] - Set by the server on a receipt created from an email: the System Email login (IMAP username) of the inbox the email was read from. Ignored on write.
/// * [emailSender] - Set by the server on a receipt created from an email: the sender's From address. Ignored on write.
/// * [claimId] - Set by the server: the id of the reimbursement claim this receipt is in. Read-only (ignored on write). Omitted when the receipt is in no claim or its claim is not visible to the caller.
/// * [claimName] - Set by the server: the name of the claim in claimId. Read-only (ignored on write). Omitted when the receipt is in no claim or its claim is not visible to the caller.
@BuiltValue()
abstract class Receipt implements Built<Receipt, ReceiptBuilder> {
  /// Receipt total amount
  @BuiltValueField(wireName: r'amount')
  String get amount;

  /// ISO 4217 currency code of the receipt's amount and of its items' amounts (items have no currency of their own). A plain string, never an enum. Never empty on a stored receipt.
  @BuiltValueField(wireName: r'currency')
  String? get currency;

  /// What the receipt came to in its group's base currency (decimal string), as the user recorded or the server derived it. Only present when the currency differs from the base currency and a conversion was recorded.
  @BuiltValueField(wireName: r'baseAmount')
  String? get baseAmount;

  /// Read-only, server-set. The ISO 4217 currency baseAmount is in: the group's base currency at the moment the conversion was written. A conversion counts as converted (claim figures, packet) only while this equals the group's current base currency; otherwise the receipt is treated as unconverted and the stored conversion is shown as recorded. Not to be confused with GroupReceiptSettings.baseCurrency.
  @BuiltValueField(wireName: r'baseAmountCurrency')
  String? get baseAmountCurrency;

  /// Base-currency units per one unit of the receipt's currency (decimal string). Only present alongside baseAmount.
  @BuiltValueField(wireName: r'exchangeRate')
  String? get exchangeRate;

  /// Calendar day (yyyy-MM-dd) the exchange rate applies to
  @BuiltValueField(wireName: r'exchangeRateDate')
  String? get exchangeRateDate;

  /// Free text saying where the rate came from, e.g. \"Card statement\"
  @BuiltValueField(wireName: r'exchangeRateSource')
  String? get exchangeRateSource;

  /// Read-only warning, present when baseAmount and exchangeRate disagree about the amount by more than 1%. Never blocks a save.
  @BuiltValueField(wireName: r'conversionNote')
  String? get conversionNote;

  /// Categories associated to receipt
  @BuiltValueField(wireName: r'categories')
  BuiltList<Category> get categories;

  /// Comments associated to receipt
  @BuiltValueField(wireName: r'comments')
  BuiltList<Comment> get comments;

  /// Custom fields associated to receipt
  @BuiltValueField(wireName: r'customFields')
  BuiltList<CustomFieldValue> get customFields;

  @BuiltValueField(wireName: r'createdAt')
  String? get createdAt;

  @BuiltValueField(wireName: r'createdBy')
  int? get createdBy;

  /// Receipt date
  @BuiltValueField(wireName: r'date')
  String get date;

  /// Text of the receipt's earliest comment the caller may see. Only populated by the paged receipt list (POST /receipt/group/{groupId}); absent when there is no such comment.
  @BuiltValueField(wireName: r'firstComment')
  String? get firstComment;

  /// Group foreign key
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  @BuiltValueField(wireName: r'id')
  int get id;

  /// Files associated to receipt
  @BuiltValueField(wireName: r'imageFiles')
  BuiltList<FileData>? get imageFiles;

  /// Receipt name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// User paid foreign key
  @BuiltValueField(wireName: r'paidByUserId')
  int get paidByUserId;

  /// Items associated to receipt
  @BuiltValueField(wireName: r'receiptItems')
  BuiltList<Item> get receiptItems;

  /// Date resolved
  @BuiltValueField(wireName: r'resolvedDate')
  String? get resolvedDate;

  @BuiltValueField(wireName: r'status')
  ReceiptStatus get status;
  // enum statusEnum {  OPEN,  NEEDS_ATTENTION,  RESOLVED,  DRAFT,  DECLINED,  ,  };

  /// Tags associated to receipt
  @BuiltValueField(wireName: r'tags')
  BuiltList<Tag> get tags;

  @BuiltValueField(wireName: r'updatedAt')
  String? get updatedAt;

  /// Created by string, which is anything that is not a user
  @BuiltValueField(wireName: r'createdByString')
  String? get createdByString;

  /// Set by the server on a receipt created from an email: the System Email login (IMAP username) of the inbox the email was read from. Ignored on write.
  @BuiltValueField(wireName: r'emailInbox')
  String? get emailInbox;

  /// Set by the server on a receipt created from an email: the sender's From address. Ignored on write.
  @BuiltValueField(wireName: r'emailSender')
  String? get emailSender;

  /// Set by the server: the id of the reimbursement claim this receipt is in. Read-only (ignored on write). Omitted when the receipt is in no claim or its claim is not visible to the caller.
  @BuiltValueField(wireName: r'claimId')
  int? get claimId;

  /// Set by the server: the name of the claim in claimId. Read-only (ignored on write). Omitted when the receipt is in no claim or its claim is not visible to the caller.
  @BuiltValueField(wireName: r'claimName')
  String? get claimName;

  Receipt._();

  factory Receipt([void updates(ReceiptBuilder b)]) = _$Receipt;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(ReceiptBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<Receipt> get serializer => _$ReceiptSerializer();
}

class _$ReceiptSerializer implements PrimitiveSerializer<Receipt> {
  @override
  final Iterable<Type> types = const [Receipt, _$Receipt];

  @override
  final String wireName = r'Receipt';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    Receipt object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'amount';
    yield serializers.serialize(
      object.amount,
      specifiedType: const FullType(String),
    );
    if (object.currency != null) {
      yield r'currency';
      yield serializers.serialize(
        object.currency,
        specifiedType: const FullType(String),
      );
    }
    if (object.baseAmount != null) {
      yield r'baseAmount';
      yield serializers.serialize(
        object.baseAmount,
        specifiedType: const FullType(String),
      );
    }
    if (object.baseAmountCurrency != null) {
      yield r'baseAmountCurrency';
      yield serializers.serialize(
        object.baseAmountCurrency,
        specifiedType: const FullType(String),
      );
    }
    if (object.exchangeRate != null) {
      yield r'exchangeRate';
      yield serializers.serialize(
        object.exchangeRate,
        specifiedType: const FullType(String),
      );
    }
    if (object.exchangeRateDate != null) {
      yield r'exchangeRateDate';
      yield serializers.serialize(
        object.exchangeRateDate,
        specifiedType: const FullType(String),
      );
    }
    if (object.exchangeRateSource != null) {
      yield r'exchangeRateSource';
      yield serializers.serialize(
        object.exchangeRateSource,
        specifiedType: const FullType(String),
      );
    }
    if (object.conversionNote != null) {
      yield r'conversionNote';
      yield serializers.serialize(
        object.conversionNote,
        specifiedType: const FullType(String),
      );
    }
    yield r'categories';
    yield serializers.serialize(
      object.categories,
      specifiedType: const FullType(BuiltList, [FullType(Category)]),
    );
    yield r'comments';
    yield serializers.serialize(
      object.comments,
      specifiedType: const FullType(BuiltList, [FullType(Comment)]),
    );
    yield r'customFields';
    yield serializers.serialize(
      object.customFields,
      specifiedType: const FullType(BuiltList, [FullType(CustomFieldValue)]),
    );
    if (object.createdAt != null) {
      yield r'createdAt';
      yield serializers.serialize(
        object.createdAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdBy != null) {
      yield r'createdBy';
      yield serializers.serialize(
        object.createdBy,
        specifiedType: const FullType(int),
      );
    }
    yield r'date';
    yield serializers.serialize(
      object.date,
      specifiedType: const FullType(String),
    );
    if (object.firstComment != null) {
      yield r'firstComment';
      yield serializers.serialize(
        object.firstComment,
        specifiedType: const FullType(String),
      );
    }
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    if (object.imageFiles != null) {
      yield r'imageFiles';
      yield serializers.serialize(
        object.imageFiles,
        specifiedType: const FullType(BuiltList, [FullType(FileData)]),
      );
    }
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'paidByUserId';
    yield serializers.serialize(
      object.paidByUserId,
      specifiedType: const FullType(int),
    );
    yield r'receiptItems';
    yield serializers.serialize(
      object.receiptItems,
      specifiedType: const FullType(BuiltList, [FullType(Item)]),
    );
    if (object.resolvedDate != null) {
      yield r'resolvedDate';
      yield serializers.serialize(
        object.resolvedDate,
        specifiedType: const FullType(String),
      );
    }
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(ReceiptStatus),
    );
    yield r'tags';
    yield serializers.serialize(
      object.tags,
      specifiedType: const FullType(BuiltList, [FullType(Tag)]),
    );
    if (object.updatedAt != null) {
      yield r'updatedAt';
      yield serializers.serialize(
        object.updatedAt,
        specifiedType: const FullType(String),
      );
    }
    if (object.createdByString != null) {
      yield r'createdByString';
      yield serializers.serialize(
        object.createdByString,
        specifiedType: const FullType(String),
      );
    }
    if (object.emailInbox != null) {
      yield r'emailInbox';
      yield serializers.serialize(
        object.emailInbox,
        specifiedType: const FullType(String),
      );
    }
    if (object.emailSender != null) {
      yield r'emailSender';
      yield serializers.serialize(
        object.emailSender,
        specifiedType: const FullType(String),
      );
    }
    if (object.claimId != null) {
      yield r'claimId';
      yield serializers.serialize(
        object.claimId,
        specifiedType: const FullType(int),
      );
    }
    if (object.claimName != null) {
      yield r'claimName';
      yield serializers.serialize(
        object.claimName,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    Receipt object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required ReceiptBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'amount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.amount = valueDes;
          break;
        case r'currency':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.currency = valueDes;
          break;
        case r'baseAmount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.baseAmount = valueDes;
          break;
        case r'baseAmountCurrency':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.baseAmountCurrency = valueDes;
          break;
        case r'exchangeRate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.exchangeRate = valueDes;
          break;
        case r'exchangeRateDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.exchangeRateDate = valueDes;
          break;
        case r'exchangeRateSource':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.exchangeRateSource = valueDes;
          break;
        case r'conversionNote':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.conversionNote = valueDes;
          break;
        case r'categories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Category)]),
          ) as BuiltList<Category>;
          result.categories.replace(valueDes);
          break;
        case r'comments':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Comment)]),
          ) as BuiltList<Comment>;
          result.comments.replace(valueDes);
          break;
        case r'customFields':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(CustomFieldValue)]),
          ) as BuiltList<CustomFieldValue>;
          result.customFields.replace(valueDes);
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdAt = valueDes;
          break;
        case r'createdBy':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdBy = valueDes;
          break;
        case r'date':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.date = valueDes;
          break;
        case r'firstComment':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.firstComment = valueDes;
          break;
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'imageFiles':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(FileData)]),
          ) as BuiltList<FileData>;
          result.imageFiles.replace(valueDes);
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'paidByUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.paidByUserId = valueDes;
          break;
        case r'receiptItems':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Item)]),
          ) as BuiltList<Item>;
          result.receiptItems.replace(valueDes);
          break;
        case r'resolvedDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.resolvedDate = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptStatus),
          ) as ReceiptStatus;
          result.status = valueDes;
          break;
        case r'tags':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(Tag)]),
          ) as BuiltList<Tag>;
          result.tags.replace(valueDes);
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.updatedAt = valueDes;
          break;
        case r'createdByString':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.createdByString = valueDes;
          break;
        case r'emailInbox':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.emailInbox = valueDes;
          break;
        case r'emailSender':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.emailSender = valueDes;
          break;
        case r'claimId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.claimId = valueDes;
          break;
        case r'claimName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.claimName = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  Receipt deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = ReceiptBuilder();
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

