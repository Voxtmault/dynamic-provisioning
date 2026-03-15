import 'package:flutter/foundation.dart';

import '../models/tenant.dart';
import '../services/api_service.dart';

enum TenantState { idle, loading, loaded, error }

class TenantProvider extends ChangeNotifier {
  final ApiService _api;

  TenantProvider(this._api);

  TenantState _state = TenantState.idle;
  List<TenantListItem> _tenants = [];
  String _errorMessage = '';

  TenantState get state => _state;
  List<TenantListItem> get tenants => _tenants;
  String get errorMessage => _errorMessage;

  Future<void> loadTenants(String token) async {
    _state = TenantState.loading;
    _errorMessage = '';
    notifyListeners();

    try {
      _tenants = await _api.listTenants(token);
      _state = TenantState.loaded;
    } catch (e) {
      _state = TenantState.error;
      _errorMessage = _extractMessage(e);
    }

    notifyListeners();
  }

  /// Registers a new tenant, optionally uploading a photo to the presigned URL.
  /// Returns the [RegisterTenantResponse] so the caller can navigate to the
  /// success screen.
  Future<RegisterTenantResponse> registerTenant({
    required String token,
    required String name,
    required List<String> colorPalette,
    Uint8List? photoBytes,
    String? photoMimeType,
  }) async {
    final resp = await _api.registerTenant(token, name, colorPalette);

    if (photoBytes != null &&
        photoBytes.isNotEmpty &&
        resp.photoUploadUrl.isNotEmpty) {
      await _api.uploadPhoto(
        resp.photoUploadUrl,
        photoBytes,
        photoMimeType ?? 'image/jpeg',
      );
    }

    // Reload list in the background so dashboard stays fresh.
    loadTenants(token);

    return resp;
  }

  static String _extractMessage(Object e) {
    final msg = e.toString();
    if (msg.startsWith('Exception: ')) return msg.substring(11);
    return msg;
  }
}
