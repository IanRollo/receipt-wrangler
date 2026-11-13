// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_summary_currency_total.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptSummaryCurrencyTotal extends ReceiptSummaryCurrencyTotal {
  @override
  final String currency;
  @override
  final int receiptCount;
  @override
  final String total;
  @override
  final BuiltList<ReceiptSummaryCustomFieldTotal> customFieldTotals;

  factory _$ReceiptSummaryCurrencyTotal(
          [void Function(ReceiptSummaryCurrencyTotalBuilder)? updates]) =>
      (ReceiptSummaryCurrencyTotalBuilder()..update(updates))._build();

  _$ReceiptSummaryCurrencyTotal._(
      {required this.currency,
      required this.receiptCount,
      required this.total,
      required this.customFieldTotals})
      : super._();
  @override
  ReceiptSummaryCurrencyTotal rebuild(
          void Function(ReceiptSummaryCurrencyTotalBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptSummaryCurrencyTotalBuilder toBuilder() =>
      ReceiptSummaryCurrencyTotalBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptSummaryCurrencyTotal &&
        currency == other.currency &&
        receiptCount == other.receiptCount &&
        total == other.total &&
        customFieldTotals == other.customFieldTotals;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, currency.hashCode);
    _$hash = $jc(_$hash, receiptCount.hashCode);
    _$hash = $jc(_$hash, total.hashCode);
    _$hash = $jc(_$hash, customFieldTotals.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptSummaryCurrencyTotal')
          ..add('currency', currency)
          ..add('receiptCount', receiptCount)
          ..add('total', total)
          ..add('customFieldTotals', customFieldTotals))
        .toString();
  }
}

class ReceiptSummaryCurrencyTotalBuilder
    implements
        Builder<ReceiptSummaryCurrencyTotal,
            ReceiptSummaryCurrencyTotalBuilder> {
  _$ReceiptSummaryCurrencyTotal? _$v;

  String? _currency;
  String? get currency => _$this._currency;
  set currency(String? currency) => _$this._currency = currency;

  int? _receiptCount;
  int? get receiptCount => _$this._receiptCount;
  set receiptCount(int? receiptCount) => _$this._receiptCount = receiptCount;

  String? _total;
  String? get total => _$this._total;
  set total(String? total) => _$this._total = total;

  ListBuilder<ReceiptSummaryCustomFieldTotal>? _customFieldTotals;
  ListBuilder<ReceiptSummaryCustomFieldTotal> get customFieldTotals =>
      _$this._customFieldTotals ??=
          ListBuilder<ReceiptSummaryCustomFieldTotal>();
  set customFieldTotals(
          ListBuilder<ReceiptSummaryCustomFieldTotal>? customFieldTotals) =>
      _$this._customFieldTotals = customFieldTotals;

  ReceiptSummaryCurrencyTotalBuilder() {
    ReceiptSummaryCurrencyTotal._defaults(this);
  }

  ReceiptSummaryCurrencyTotalBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _currency = $v.currency;
      _receiptCount = $v.receiptCount;
      _total = $v.total;
      _customFieldTotals = $v.customFieldTotals.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptSummaryCurrencyTotal other) {
    _$v = other as _$ReceiptSummaryCurrencyTotal;
  }

  @override
  void update(void Function(ReceiptSummaryCurrencyTotalBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptSummaryCurrencyTotal build() => _build();

  _$ReceiptSummaryCurrencyTotal _build() {
    _$ReceiptSummaryCurrencyTotal _$result;
    try {
      _$result = _$v ??
          _$ReceiptSummaryCurrencyTotal._(
            currency: BuiltValueNullFieldError.checkNotNull(
                currency, r'ReceiptSummaryCurrencyTotal', 'currency'),
            receiptCount: BuiltValueNullFieldError.checkNotNull(
                receiptCount, r'ReceiptSummaryCurrencyTotal', 'receiptCount'),
            total: BuiltValueNullFieldError.checkNotNull(
                total, r'ReceiptSummaryCurrencyTotal', 'total'),
            customFieldTotals: customFieldTotals.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'customFieldTotals';
        customFieldTotals.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'ReceiptSummaryCurrencyTotal', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
