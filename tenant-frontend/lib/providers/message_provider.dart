import 'dart:async';

import 'package:flutter/foundation.dart';

import '../models/message.dart';
import '../services/api_service.dart';

enum MessageState { idle, loading, loaded, error }

class MessageProvider extends ChangeNotifier {
  final ApiService _api;
  static const Duration _pollInterval = Duration(seconds: 10);
  static const int _pageLimit = 20;

  MessageProvider(this._api);

  MessageState _state = MessageState.idle;
  List<Message> _messages = [];
  int _total = 0;
  String _errorMessage = '';
  Timer? _timer;

  MessageState get state => _state;
  List<Message> get messages => _messages;
  int get total => _total;
  String get errorMessage => _errorMessage;

  void startPolling() {
    fetchMessages();
    _timer?.cancel();
    _timer = Timer.periodic(_pollInterval, (_) => fetchMessages(silent: true));
  }

  Future<void> fetchMessages({bool silent = false}) async {
    if (!silent) {
      _state = MessageState.loading;
      notifyListeners();
    }

    try {
      final result = await _api.fetchMessages(page: 1, limit: _pageLimit);
      _messages = result.messages;
      _total = result.total;
      _state = MessageState.loaded;
    } catch (e) {
      if (!silent) {
        _state = MessageState.error;
        _errorMessage = e.toString();
      }
      // Silent poll failures are swallowed — the last good data stays visible.
    }

    notifyListeners();
  }

  Future<void> postMessage(String handlerName, String content) async {
    await _api.postMessage(handlerName, content);
    // Refresh immediately after posting so the new message appears.
    await fetchMessages(silent: true);
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }
}
