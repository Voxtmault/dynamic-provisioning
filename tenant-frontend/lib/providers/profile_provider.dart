import 'package:flutter/material.dart';

import '../models/app_profile.dart';
import '../services/api_service.dart';
import '../utils/title_setter.dart';

enum ProfileState { idle, loading, loaded, error }

class ProfileProvider extends ChangeNotifier {
  final ApiService _api;

  ProfileProvider(this._api);

  ProfileState _state = ProfileState.idle;
  AppProfile? _profile;
  String _errorMessage = '';

  ProfileState get state => _state;
  AppProfile? get profile => _profile;
  String get errorMessage => _errorMessage;
  bool get isLoaded => _state == ProfileState.loaded;

  ThemeData get themeData {
    final p = _profile;
    if (p == null) return ThemeData.light(useMaterial3: true);

    final primary = p.primaryColor;
    final secondary = p.secondaryColor;
    final background = p.backgroundColor;

    final onPrimary = _contrastColor(primary);
    final onSecondary = _contrastColor(secondary);
    final onSurface = _contrastColor(background);
    final brightness = _isLight(background) ? Brightness.light : Brightness.dark;

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      colorScheme: ColorScheme(
        brightness: brightness,
        primary: primary,
        onPrimary: onPrimary,
        secondary: secondary,
        onSecondary: onSecondary,
        error: const Color(0xFFCF6679),
        onError: Colors.white,
        surface: background,
        onSurface: onSurface,
      ),
      scaffoldBackgroundColor: background,
      appBarTheme: AppBarTheme(
        backgroundColor: primary,
        foregroundColor: onPrimary,
        elevation: 2,
      ),
      floatingActionButtonTheme: FloatingActionButtonThemeData(
        backgroundColor: primary,
        foregroundColor: onPrimary,
      ),
      cardTheme: CardThemeData(
        elevation: 1,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
    );
  }

  Future<void> loadProfile() async {
    _state = ProfileState.loading;
    _errorMessage = '';
    notifyListeners();

    try {
      _profile = await _api.fetchProfile();
      _state = ProfileState.loaded;
      setDocumentTitle(_profile!.appName);
    } catch (e) {
      _state = ProfileState.error;
      _errorMessage = e.toString();
    }

    notifyListeners();
  }

  // WCAG threshold: luminance > 0.179 means black text has better contrast.
  static bool _isLight(Color color) => color.computeLuminance() > 0.179;

  static Color _contrastColor(Color bg) =>
      _isLight(bg) ? Colors.black87 : Colors.white;
}
