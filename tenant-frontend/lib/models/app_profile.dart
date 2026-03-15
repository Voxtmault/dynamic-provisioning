import 'package:flutter/material.dart';

class AppProfile {
  final String tenantId;
  final String appName;
  final String appPhotoUrl;
  final List<String> colorPalette;

  const AppProfile({
    required this.tenantId,
    required this.appName,
    required this.appPhotoUrl,
    required this.colorPalette,
  });

  factory AppProfile.fromJson(Map<String, dynamic> json) {
    return AppProfile(
      tenantId: json['tenant_id'] as String? ?? '',
      appName: json['app_name'] as String? ?? 'Messaging Board',
      appPhotoUrl: json['app_photo_url'] as String? ?? '',
      colorPalette: (json['color_palette'] as List<dynamic>?)
              ?.map((e) => e as String)
              .toList() ??
          [],
    );
  }

  // Palette index conventions: 0=primary, 1=secondary, 2=background
  Color get primaryColor => _hexToColor(colorPalette.elementAtOrNull(0), const Color(0xFF6200EE));
  Color get secondaryColor => _hexToColor(colorPalette.elementAtOrNull(1), const Color(0xFF03DAC6));
  Color get backgroundColor => _hexToColor(colorPalette.elementAtOrNull(2), const Color(0xFFF5F5F5));

  static Color _hexToColor(String? hex, Color fallback) {
    if (hex == null || hex.isEmpty) return fallback;
    final cleaned = hex.replaceFirst('#', '');
    final value = int.tryParse(
      cleaned.length == 6 ? 'FF$cleaned' : cleaned,
      radix: 16,
    );
    return value != null ? Color(value) : fallback;
  }
}
