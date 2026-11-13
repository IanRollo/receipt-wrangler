// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'currency_info.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CurrencyInfo extends CurrencyInfo {
  @override
  final String code;
  @override
  final int exponent;
  @override
  final String? symbol;

  factory _$CurrencyInfo([void Function(CurrencyInfoBuilder)? updates]) =>
      (CurrencyInfoBuilder()..update(updates))._build();

  _$CurrencyInfo._({required this.code, required this.exponent, this.symbol})
      : super._();
  @override
  CurrencyInfo rebuild(void Function(CurrencyInfoBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CurrencyInfoBuilder toBuilder() => CurrencyInfoBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CurrencyInfo &&
        code == other.code &&
        exponent == other.exponent &&
        symbol == other.symbol;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, code.hashCode);
    _$hash = $jc(_$hash, exponent.hashCode);
    _$hash = $jc(_$hash, symbol.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'CurrencyInfo')
          ..add('code', code)
          ..add('exponent', exponent)
          ..add('symbol', symbol))
        .toString();
  }
}

class CurrencyInfoBuilder
    implements Builder<CurrencyInfo, CurrencyInfoBuilder> {
  _$CurrencyInfo? _$v;

  String? _code;
  String? get code => _$this._code;
  set code(String? code) => _$this._code = code;

  int? _exponent;
  int? get exponent => _$this._exponent;
  set exponent(int? exponent) => _$this._exponent = exponent;

  String? _symbol;
  String? get symbol => _$this._symbol;
  set symbol(String? symbol) => _$this._symbol = symbol;

  CurrencyInfoBuilder() {
    CurrencyInfo._defaults(this);
  }

  CurrencyInfoBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _code = $v.code;
      _exponent = $v.exponent;
      _symbol = $v.symbol;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(CurrencyInfo other) {
    _$v = other as _$CurrencyInfo;
  }

  @override
  void update(void Function(CurrencyInfoBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CurrencyInfo build() => _build();

  _$CurrencyInfo _build() {
    final _$result = _$v ??
        _$CurrencyInfo._(
          code: BuiltValueNullFieldError.checkNotNull(
              code, r'CurrencyInfo', 'code'),
          exponent: BuiltValueNullFieldError.checkNotNull(
              exponent, r'CurrencyInfo', 'exponent'),
          symbol: symbol,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
