import 'dart:convert';
import 'dart:typed_data';

import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:http/http.dart' as http;

import '../models/tenant.dart';

class ApiService {
  // On web: turns e.g. "https://admin.example.com" → "https://api.admin.example.com/api"
  // Native fallback for unit tests / non-web debug builds.
  static String get _baseUrl {
    if (kIsWeb) {
      final base = Uri.base;
      return '${base.scheme}://api.${base.host}/api';
    }
    return 'http://localhost:8080/api';
  }

  final http.Client _client;

  ApiService({http.Client? client}) : _client = client ?? http.Client();

  // ── Auth ────────────────────────────────────────────────────────────────────

  Future<String> login(String email, String password) async {
    final uri = Uri.parse('$_baseUrl/auth/login');
    final response = await _client.post(
      uri,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );
    _assertOk(response, 'login');
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return data['token'] as String;
  }

  // ── Tenants ─────────────────────────────────────────────────────────────────

  Future<List<TenantListItem>> listTenants(String token) async {
    final uri = Uri.parse('$_baseUrl/tenants');
    final response =
        await _client.get(uri, headers: _authHeaders(token));
    _assertOk(response, 'listTenants');
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    final rawList = body['data'] as List<dynamic>? ?? [];
    return rawList
        .map((e) => TenantListItem.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<RegisterTenantResponse> registerTenant(
    String token,
    String name,
    List<String> colorPalette,
  ) async {
    final uri = Uri.parse('$_baseUrl/tenants');
    final response = await _client.post(
      uri,
      headers: {
        'Content-Type': 'application/json',
        ..._authHeaders(token),
      },
      body: jsonEncode({
        'name': name,
        'color_palette': colorPalette,
      }),
    );
    _assertOk(response, 'registerTenant');
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    return RegisterTenantResponse.fromJson(
        body['data'] as Map<String, dynamic>);
  }

  // ── Photo upload ─────────────────────────────────────────────────────────────
  // PUT directly to the presigned S3 URL — no Authorization header (the
  // presigned URL itself encodes the credentials).

  Future<void> uploadPhoto(
    String presignedUrl,
    Uint8List bytes,
    String mimeType,
  ) async {
    final uri = Uri.parse(presignedUrl);
    final response = await _client.put(
      uri,
      headers: {'Content-Type': mimeType},
      body: bytes,
    );
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw Exception(
          'Photo upload failed (HTTP ${response.statusCode}): ${response.body}');
    }
  }

  // ── Helpers ──────────────────────────────────────────────────────────────────

  Map<String, String> _authHeaders(String token) =>
      {'Authorization': 'Bearer $token'};

  void _assertOk(http.Response response, String op) {
    if (response.statusCode < 200 || response.statusCode >= 300) {
      String message = '';
      try {
        final body = jsonDecode(response.body) as Map<String, dynamic>;
        message = body['message'] as String? ?? response.body;
      } catch (_) {
        message = response.body;
      }
      throw Exception('$op failed (HTTP ${response.statusCode}): $message');
    }
  }
}
