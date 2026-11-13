// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'system_settings.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SystemSettings extends SystemSettings {
  @override
  final bool? mcpEnabled;
  @override
  final int? pdfDpi;
  @override
  final int? tempFileRetentionHours;
  @override
  final String? currencyDisplay;
  @override
  final bool? currencyHideDecimalPlaces;
  @override
  final CurrencySeparator? currencyDecimalSeparator;
  @override
  final String? mobileServerUrl;
  @override
  final String? defaultCurrency;
  @override
  final int? numWorkers;
  @override
  final bool? enableLocalSignUp;
  @override
  final int? mcpRefreshTokenValidForHours;
  @override
  final CurrencySeparator? currencyThousandthsSeparator;
  @override
  final bool? secureCookies;
  @override
  final bool? showLoginQr;
  @override
  final String? serverPublicUrl;
  @override
  final String? timeZone;
  @override
  final String? mcpPublicUrl;
  @override
  final bool? debugOcr;
  @override
  final int? fallbackReceiptProcessingSettingsId;
  @override
  final int? receiptProcessingSettingsId;
  @override
  final int? refreshTokenValidForHours;
  @override
  final CurrencySymbolPosition? currencySymbolPosition;
  @override
  final int? taskConcurrency;
  @override
  final int? emailPollingInterval;
  @override
  final BuiltList<TaskQueueConfiguration> taskQueueConfigurations;
  @override
  final int id;
  @override
  final String createdAt;
  @override
  final int? createdBy;
  @override
  final String? createdByString;
  @override
  final String? updatedAt;

  factory _$SystemSettings([void Function(SystemSettingsBuilder)? updates]) =>
      (SystemSettingsBuilder()..update(updates))._build();

  _$SystemSettings._(
      {this.mcpEnabled,
      this.pdfDpi,
      this.tempFileRetentionHours,
      this.currencyDisplay,
      this.currencyHideDecimalPlaces,
      this.currencyDecimalSeparator,
      this.mobileServerUrl,
      this.defaultCurrency,
      this.numWorkers,
      this.enableLocalSignUp,
      this.mcpRefreshTokenValidForHours,
      this.currencyThousandthsSeparator,
      this.secureCookies,
      this.showLoginQr,
      this.serverPublicUrl,
      this.timeZone,
      this.mcpPublicUrl,
      this.debugOcr,
      this.fallbackReceiptProcessingSettingsId,
      this.receiptProcessingSettingsId,
      this.refreshTokenValidForHours,
      this.currencySymbolPosition,
      this.taskConcurrency,
      this.emailPollingInterval,
      required this.taskQueueConfigurations,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._();
  @override
  SystemSettings rebuild(void Function(SystemSettingsBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SystemSettingsBuilder toBuilder() => SystemSettingsBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SystemSettings &&
        mcpEnabled == other.mcpEnabled &&
        pdfDpi == other.pdfDpi &&
        tempFileRetentionHours == other.tempFileRetentionHours &&
        currencyDisplay == other.currencyDisplay &&
        currencyHideDecimalPlaces == other.currencyHideDecimalPlaces &&
        currencyDecimalSeparator == other.currencyDecimalSeparator &&
        mobileServerUrl == other.mobileServerUrl &&
        defaultCurrency == other.defaultCurrency &&
        numWorkers == other.numWorkers &&
        enableLocalSignUp == other.enableLocalSignUp &&
        mcpRefreshTokenValidForHours == other.mcpRefreshTokenValidForHours &&
        currencyThousandthsSeparator == other.currencyThousandthsSeparator &&
        secureCookies == other.secureCookies &&
        showLoginQr == other.showLoginQr &&
        serverPublicUrl == other.serverPublicUrl &&
        timeZone == other.timeZone &&
        mcpPublicUrl == other.mcpPublicUrl &&
        debugOcr == other.debugOcr &&
        fallbackReceiptProcessingSettingsId ==
            other.fallbackReceiptProcessingSettingsId &&
        receiptProcessingSettingsId == other.receiptProcessingSettingsId &&
        refreshTokenValidForHours == other.refreshTokenValidForHours &&
        currencySymbolPosition == other.currencySymbolPosition &&
        taskConcurrency == other.taskConcurrency &&
        emailPollingInterval == other.emailPollingInterval &&
        taskQueueConfigurations == other.taskQueueConfigurations &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, mcpEnabled.hashCode);
    _$hash = $jc(_$hash, pdfDpi.hashCode);
    _$hash = $jc(_$hash, tempFileRetentionHours.hashCode);
    _$hash = $jc(_$hash, currencyDisplay.hashCode);
    _$hash = $jc(_$hash, currencyHideDecimalPlaces.hashCode);
    _$hash = $jc(_$hash, currencyDecimalSeparator.hashCode);
    _$hash = $jc(_$hash, mobileServerUrl.hashCode);
    _$hash = $jc(_$hash, defaultCurrency.hashCode);
    _$hash = $jc(_$hash, numWorkers.hashCode);
    _$hash = $jc(_$hash, enableLocalSignUp.hashCode);
    _$hash = $jc(_$hash, mcpRefreshTokenValidForHours.hashCode);
    _$hash = $jc(_$hash, currencyThousandthsSeparator.hashCode);
    _$hash = $jc(_$hash, secureCookies.hashCode);
    _$hash = $jc(_$hash, showLoginQr.hashCode);
    _$hash = $jc(_$hash, serverPublicUrl.hashCode);
    _$hash = $jc(_$hash, timeZone.hashCode);
    _$hash = $jc(_$hash, mcpPublicUrl.hashCode);
    _$hash = $jc(_$hash, debugOcr.hashCode);
    _$hash = $jc(_$hash, fallbackReceiptProcessingSettingsId.hashCode);
    _$hash = $jc(_$hash, receiptProcessingSettingsId.hashCode);
    _$hash = $jc(_$hash, refreshTokenValidForHours.hashCode);
    _$hash = $jc(_$hash, currencySymbolPosition.hashCode);
    _$hash = $jc(_$hash, taskConcurrency.hashCode);
    _$hash = $jc(_$hash, emailPollingInterval.hashCode);
    _$hash = $jc(_$hash, taskQueueConfigurations.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SystemSettings')
          ..add('mcpEnabled', mcpEnabled)
          ..add('pdfDpi', pdfDpi)
          ..add('tempFileRetentionHours', tempFileRetentionHours)
          ..add('currencyDisplay', currencyDisplay)
          ..add('currencyHideDecimalPlaces', currencyHideDecimalPlaces)
          ..add('currencyDecimalSeparator', currencyDecimalSeparator)
          ..add('mobileServerUrl', mobileServerUrl)
          ..add('defaultCurrency', defaultCurrency)
          ..add('numWorkers', numWorkers)
          ..add('enableLocalSignUp', enableLocalSignUp)
          ..add('mcpRefreshTokenValidForHours', mcpRefreshTokenValidForHours)
          ..add('currencyThousandthsSeparator', currencyThousandthsSeparator)
          ..add('secureCookies', secureCookies)
          ..add('showLoginQr', showLoginQr)
          ..add('serverPublicUrl', serverPublicUrl)
          ..add('timeZone', timeZone)
          ..add('mcpPublicUrl', mcpPublicUrl)
          ..add('debugOcr', debugOcr)
          ..add('fallbackReceiptProcessingSettingsId',
              fallbackReceiptProcessingSettingsId)
          ..add('receiptProcessingSettingsId', receiptProcessingSettingsId)
          ..add('refreshTokenValidForHours', refreshTokenValidForHours)
          ..add('currencySymbolPosition', currencySymbolPosition)
          ..add('taskConcurrency', taskConcurrency)
          ..add('emailPollingInterval', emailPollingInterval)
          ..add('taskQueueConfigurations', taskQueueConfigurations)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class SystemSettingsBuilder
    implements
        Builder<SystemSettings, SystemSettingsBuilder>,
        BaseModelBuilder {
  _$SystemSettings? _$v;

  bool? _mcpEnabled;
  bool? get mcpEnabled => _$this._mcpEnabled;
  set mcpEnabled(covariant bool? mcpEnabled) => _$this._mcpEnabled = mcpEnabled;

  int? _pdfDpi;
  int? get pdfDpi => _$this._pdfDpi;
  set pdfDpi(covariant int? pdfDpi) => _$this._pdfDpi = pdfDpi;

  int? _tempFileRetentionHours;
  int? get tempFileRetentionHours => _$this._tempFileRetentionHours;
  set tempFileRetentionHours(covariant int? tempFileRetentionHours) =>
      _$this._tempFileRetentionHours = tempFileRetentionHours;

  String? _currencyDisplay;
  String? get currencyDisplay => _$this._currencyDisplay;
  set currencyDisplay(covariant String? currencyDisplay) =>
      _$this._currencyDisplay = currencyDisplay;

  bool? _currencyHideDecimalPlaces;
  bool? get currencyHideDecimalPlaces => _$this._currencyHideDecimalPlaces;
  set currencyHideDecimalPlaces(covariant bool? currencyHideDecimalPlaces) =>
      _$this._currencyHideDecimalPlaces = currencyHideDecimalPlaces;

  CurrencySeparator? _currencyDecimalSeparator;
  CurrencySeparator? get currencyDecimalSeparator =>
      _$this._currencyDecimalSeparator;
  set currencyDecimalSeparator(
          covariant CurrencySeparator? currencyDecimalSeparator) =>
      _$this._currencyDecimalSeparator = currencyDecimalSeparator;

  String? _mobileServerUrl;
  String? get mobileServerUrl => _$this._mobileServerUrl;
  set mobileServerUrl(covariant String? mobileServerUrl) =>
      _$this._mobileServerUrl = mobileServerUrl;

  String? _defaultCurrency;
  String? get defaultCurrency => _$this._defaultCurrency;
  set defaultCurrency(covariant String? defaultCurrency) =>
      _$this._defaultCurrency = defaultCurrency;

  int? _numWorkers;
  int? get numWorkers => _$this._numWorkers;
  set numWorkers(covariant int? numWorkers) => _$this._numWorkers = numWorkers;

  bool? _enableLocalSignUp;
  bool? get enableLocalSignUp => _$this._enableLocalSignUp;
  set enableLocalSignUp(covariant bool? enableLocalSignUp) =>
      _$this._enableLocalSignUp = enableLocalSignUp;

  int? _mcpRefreshTokenValidForHours;
  int? get mcpRefreshTokenValidForHours => _$this._mcpRefreshTokenValidForHours;
  set mcpRefreshTokenValidForHours(
          covariant int? mcpRefreshTokenValidForHours) =>
      _$this._mcpRefreshTokenValidForHours = mcpRefreshTokenValidForHours;

  CurrencySeparator? _currencyThousandthsSeparator;
  CurrencySeparator? get currencyThousandthsSeparator =>
      _$this._currencyThousandthsSeparator;
  set currencyThousandthsSeparator(
          covariant CurrencySeparator? currencyThousandthsSeparator) =>
      _$this._currencyThousandthsSeparator = currencyThousandthsSeparator;

  bool? _secureCookies;
  bool? get secureCookies => _$this._secureCookies;
  set secureCookies(covariant bool? secureCookies) =>
      _$this._secureCookies = secureCookies;

  bool? _showLoginQr;
  bool? get showLoginQr => _$this._showLoginQr;
  set showLoginQr(covariant bool? showLoginQr) =>
      _$this._showLoginQr = showLoginQr;

  String? _serverPublicUrl;
  String? get serverPublicUrl => _$this._serverPublicUrl;
  set serverPublicUrl(covariant String? serverPublicUrl) =>
      _$this._serverPublicUrl = serverPublicUrl;

  String? _timeZone;
  String? get timeZone => _$this._timeZone;
  set timeZone(covariant String? timeZone) => _$this._timeZone = timeZone;

  String? _mcpPublicUrl;
  String? get mcpPublicUrl => _$this._mcpPublicUrl;
  set mcpPublicUrl(covariant String? mcpPublicUrl) =>
      _$this._mcpPublicUrl = mcpPublicUrl;

  bool? _debugOcr;
  bool? get debugOcr => _$this._debugOcr;
  set debugOcr(covariant bool? debugOcr) => _$this._debugOcr = debugOcr;

  int? _fallbackReceiptProcessingSettingsId;
  int? get fallbackReceiptProcessingSettingsId =>
      _$this._fallbackReceiptProcessingSettingsId;
  set fallbackReceiptProcessingSettingsId(
          covariant int? fallbackReceiptProcessingSettingsId) =>
      _$this._fallbackReceiptProcessingSettingsId =
          fallbackReceiptProcessingSettingsId;

  int? _receiptProcessingSettingsId;
  int? get receiptProcessingSettingsId => _$this._receiptProcessingSettingsId;
  set receiptProcessingSettingsId(covariant int? receiptProcessingSettingsId) =>
      _$this._receiptProcessingSettingsId = receiptProcessingSettingsId;

  int? _refreshTokenValidForHours;
  int? get refreshTokenValidForHours => _$this._refreshTokenValidForHours;
  set refreshTokenValidForHours(covariant int? refreshTokenValidForHours) =>
      _$this._refreshTokenValidForHours = refreshTokenValidForHours;

  CurrencySymbolPosition? _currencySymbolPosition;
  CurrencySymbolPosition? get currencySymbolPosition =>
      _$this._currencySymbolPosition;
  set currencySymbolPosition(
          covariant CurrencySymbolPosition? currencySymbolPosition) =>
      _$this._currencySymbolPosition = currencySymbolPosition;

  int? _taskConcurrency;
  int? get taskConcurrency => _$this._taskConcurrency;
  set taskConcurrency(covariant int? taskConcurrency) =>
      _$this._taskConcurrency = taskConcurrency;

  int? _emailPollingInterval;
  int? get emailPollingInterval => _$this._emailPollingInterval;
  set emailPollingInterval(covariant int? emailPollingInterval) =>
      _$this._emailPollingInterval = emailPollingInterval;

  ListBuilder<TaskQueueConfiguration>? _taskQueueConfigurations;
  ListBuilder<TaskQueueConfiguration> get taskQueueConfigurations =>
      _$this._taskQueueConfigurations ??= ListBuilder<TaskQueueConfiguration>();
  set taskQueueConfigurations(
          covariant ListBuilder<TaskQueueConfiguration>?
              taskQueueConfigurations) =>
      _$this._taskQueueConfigurations = taskQueueConfigurations;

  int? _id;
  int? get id => _$this._id;
  set id(covariant int? id) => _$this._id = id;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(covariant String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(covariant int? createdBy) => _$this._createdBy = createdBy;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(covariant String? createdByString) =>
      _$this._createdByString = createdByString;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(covariant String? updatedAt) => _$this._updatedAt = updatedAt;

  SystemSettingsBuilder() {
    SystemSettings._defaults(this);
  }

  SystemSettingsBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _mcpEnabled = $v.mcpEnabled;
      _pdfDpi = $v.pdfDpi;
      _tempFileRetentionHours = $v.tempFileRetentionHours;
      _currencyDisplay = $v.currencyDisplay;
      _currencyHideDecimalPlaces = $v.currencyHideDecimalPlaces;
      _currencyDecimalSeparator = $v.currencyDecimalSeparator;
      _mobileServerUrl = $v.mobileServerUrl;
      _defaultCurrency = $v.defaultCurrency;
      _numWorkers = $v.numWorkers;
      _enableLocalSignUp = $v.enableLocalSignUp;
      _mcpRefreshTokenValidForHours = $v.mcpRefreshTokenValidForHours;
      _currencyThousandthsSeparator = $v.currencyThousandthsSeparator;
      _secureCookies = $v.secureCookies;
      _showLoginQr = $v.showLoginQr;
      _serverPublicUrl = $v.serverPublicUrl;
      _timeZone = $v.timeZone;
      _mcpPublicUrl = $v.mcpPublicUrl;
      _debugOcr = $v.debugOcr;
      _fallbackReceiptProcessingSettingsId =
          $v.fallbackReceiptProcessingSettingsId;
      _receiptProcessingSettingsId = $v.receiptProcessingSettingsId;
      _refreshTokenValidForHours = $v.refreshTokenValidForHours;
      _currencySymbolPosition = $v.currencySymbolPosition;
      _taskConcurrency = $v.taskConcurrency;
      _emailPollingInterval = $v.emailPollingInterval;
      _taskQueueConfigurations = $v.taskQueueConfigurations.toBuilder();
      _id = $v.id;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _createdByString = $v.createdByString;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant SystemSettings other) {
    _$v = other as _$SystemSettings;
  }

  @override
  void update(void Function(SystemSettingsBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SystemSettings build() => _build();

  _$SystemSettings _build() {
    _$SystemSettings _$result;
    try {
      _$result = _$v ??
          _$SystemSettings._(
            mcpEnabled: mcpEnabled,
            pdfDpi: pdfDpi,
            tempFileRetentionHours: tempFileRetentionHours,
            currencyDisplay: currencyDisplay,
            currencyHideDecimalPlaces: currencyHideDecimalPlaces,
            currencyDecimalSeparator: currencyDecimalSeparator,
            mobileServerUrl: mobileServerUrl,
            defaultCurrency: defaultCurrency,
            numWorkers: numWorkers,
            enableLocalSignUp: enableLocalSignUp,
            mcpRefreshTokenValidForHours: mcpRefreshTokenValidForHours,
            currencyThousandthsSeparator: currencyThousandthsSeparator,
            secureCookies: secureCookies,
            showLoginQr: showLoginQr,
            serverPublicUrl: serverPublicUrl,
            timeZone: timeZone,
            mcpPublicUrl: mcpPublicUrl,
            debugOcr: debugOcr,
            fallbackReceiptProcessingSettingsId:
                fallbackReceiptProcessingSettingsId,
            receiptProcessingSettingsId: receiptProcessingSettingsId,
            refreshTokenValidForHours: refreshTokenValidForHours,
            currencySymbolPosition: currencySymbolPosition,
            taskConcurrency: taskConcurrency,
            emailPollingInterval: emailPollingInterval,
            taskQueueConfigurations: taskQueueConfigurations.build(),
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'SystemSettings', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'SystemSettings', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'taskQueueConfigurations';
        taskQueueConfigurations.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'SystemSettings', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
