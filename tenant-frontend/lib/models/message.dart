class Message {
  final int id;
  final String handlerName;
  final String content;
  final DateTime createdAt;

  const Message({
    required this.id,
    required this.handlerName,
    required this.content,
    required this.createdAt,
  });

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: (json['id'] as num).toInt(),
      handlerName: json['handler_name'] as String? ?? 'Anonymous',
      content: json['content'] as String? ?? '',
      // Parse as UTC then convert to local in the UI layer
      createdAt: DateTime.parse(json['created_at'] as String).toLocal(),
    );
  }
}
