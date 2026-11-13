// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim_currency_total.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ClaimCurrencyTotal extends ClaimCurrencyTotal {
  @override
  final String currency;
  @override
  final String total;

  factory _$ClaimCurrencyTotal(
          [void Function(ClaimCurrencyTotalBuilder)? updates]) =>
      (ClaimCurrencyTotalBuilder()..update(updates))._build();

  _$ClaimCurrencyTotal._({required this.currency, required this.total})
      : super._();
  @override
  ClaimCurrencyTotal rebuild(
          void Function(ClaimCurrencyTotalBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ClaimCurrencyTotalBuilder toBuilder() =>
      ClaimCurrencyTotalBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ClaimCurrencyTotal &&
        currency == other.currency &&
        total == other.total;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, currency.hashCode);
    _$hash = $jc(_$hash, total.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ClaimCurrencyTotal')
          ..add('currency', currency)
          ..add('total', total))
        .toString();
  }
}

class ClaimCurrencyTotalBuilder
    implements Builder<ClaimCurrencyTotal, ClaimCurrencyTotalBuilder> {
  _$ClaimCurrencyTotal? _$v;

  String? _currency;
  String? get currency => _$this._currency;
  set currency(String? currency) => _$this._currency = currency;

  String? _total;
  String? get total => _$this._total;
  set total(String? total) => _$this._total = total;

  ClaimCurrencyTotalBuilder() {
    ClaimCurrencyTotal._defaults(this);
  }

  ClaimCurrencyTotalBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _currency = $v.currency;
      _total = $v.total;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ClaimCurrencyTotal other) {
    _$v = other as _$ClaimCurrencyTotal;
  }

  @override
  void update(void Function(ClaimCurrencyTotalBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ClaimCurrencyTotal build() => _build();

  _$ClaimCurrencyTotal _build() {
    final _$result = _$v ??
        _$ClaimCurrencyTotal._(
          currency: BuiltValueNullFieldError.checkNotNull(
              currency, r'ClaimCurrencyTotal', 'currency'),
          total: BuiltValueNullFieldError.checkNotNull(
              total, r'ClaimCurrencyTotal', 'total'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
