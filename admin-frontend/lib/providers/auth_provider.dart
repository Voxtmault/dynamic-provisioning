import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../services/api_service.dart';

enum AuthState { unknown, unauthenticated, loading, authenticated, error }

class AuthProvider extends ChangeNotifier {
  static const _tokenKey = 'admin_jwt_token';

  final ApiService _api;

  AuthProvider(this._api);

  AuthState _state = AuthState.unknown;
  String? _token;
  String _errorMessage = '';

  AuthState get state => _state;
  String? get token => _token;
  String get errorMessage => _errorMessage;
  bool get isAuthenticated => _state == AuthState.authenticated;

  // Called once at app startup to hydrate the token from local storage.
  Future<void> init() async {
    final prefs = await SharedPreferences.getInstance();
    final stored = prefs.getString(_tokenKey);
    if (stored != null && stored.isNotEmpty) {
      _token = stored;
      _state = AuthState.authenticated;
    } else {
      _state = AuthState.unauthenticated;
    }
    notifyListeners();
  }

  Future<void> login(String email, String password) async {
    _state = AuthState.loading;
    _errorMessage = '';
    notifyListeners();

    try {
      final token = await _api.login(email, password);
      _token = token;
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString(_tokenKey, token);
      _state = AuthState.authenticated;
    } catch (e) {
      _state = AuthState.error;
      _errorMessage = _extractMessage(e);
    }

    notifyListeners();
  }

  Future<void> logout() async {
    _token = null;
    _state = AuthState.unauthenticated;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    notifyListeners();
  }

  // Extracts a human-readable message from common exception types.
  static String _extractMessage(Object e) {
    final msg = e.toString();
    // Strip the leading "Exception: " prefix Flutter adds.
    if (msg.startsWith('Exception: ')) return msg.substring(11);
    return msg;
  }
}
