import 'dart:convert';

import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:http/http.dart' as http;

import '../models/app_profile.dart';
import '../models/message.dart';

class ApiService {
  // Derives the base URL from the browser's current origin so the same
  // compiled binary works for every tenant subdomain (e.g. tenant1.domain.com).
  // Falls back to localhost for non-web builds (unit tests / native debug).
  static String get _baseUrl {
    if (kIsWeb) {
      // Uri.base is the browser's current URL; .origin gives scheme+host+port.
      final origin = Uri.base.origin;
      return '$origin/api';
    }
    return 'http://localhost:8080/api';
  }

  final http.Client _client;

  ApiService({http.Client? client}) : _client = client ?? http.Client();

  Future<AppProfile> fetchProfile() async {
    final uri = Uri.parse('$_baseUrl/profile');
    final response = await _client.get(uri, headers: _headers);
    _assertOk(response, 'fetchProfile');
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    return AppProfile.fromJson(body['data'] as Map<String, dynamic>);
  }

  Future<({List<Message> messages, int total, int page, int limit})>
      fetchMessages({int page = 1, int limit = 20}) async {
    final uri = Uri.parse('$_baseUrl/messages').replace(
      queryParameters: {'page': '$page', 'limit': '$limit'},
    );
    final response = await _client.get(uri, headers: _headers);
    _assertOk(response, 'fetchMessages');
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    final rawList = body['data'] as List<dynamic>? ?? [];
    return (
      messages: rawList
          .map((e) => Message.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: (body['total'] as num?)?.toInt() ?? 0,
      page: (body['page'] as num?)?.toInt() ?? page,
      limit: (body['limit'] as num?)?.toInt() ?? limit,
    );
  }

  Future<Message> postMessage(String handlerName, String content) async {
    final uri = Uri.parse('$_baseUrl/messages');
    final response = await _client.post(
      uri,
      headers: _headers,
      body: jsonEncode({'handler_name': handlerName, 'content': content}),
    );
    _assertOk(response, 'postMessage', expectedStatus: 201);
    final body = jsonDecode(response.body) as Map<String, dynamic>;
    return Message.fromJson(body['data'] as Map<String, dynamic>);
  }

  static const Map<String, String> _headers = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  };

  void _assertOk(http.Response response, String call, {int expectedStatus = 200}) {
    if (response.statusCode != expectedStatus) {
      throw ApiException(
        call: call,
        statusCode: response.statusCode,
        body: response.body,
      );
    }
  }

  void dispose() => _client.close();
}

class ApiException implements Exception {
  final String call;
  final int statusCode;
  final String body;

  const ApiException({
    required this.call,
    required this.statusCode,
    required this.body,
  });

  @override
  String toString() => 'ApiException[$call]: HTTP $statusCode — $body';
}
