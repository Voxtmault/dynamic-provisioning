import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/app_profile.dart';
import '../providers/message_provider.dart';
import '../providers/profile_provider.dart';
import '../widgets/error_screen.dart';
import '../widgets/loading_screen.dart';
import '../widgets/message_card.dart';
import '../widgets/post_message_dialog.dart';

class BoardScreen extends StatefulWidget {
  const BoardScreen({super.key});

  @override
  State<BoardScreen> createState() => _BoardScreenState();
}

class _BoardScreenState extends State<BoardScreen> {
  @override
  void initState() {
    super.initState();
    // Start polling once the first frame is rendered so we have a context.
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<MessageProvider>().startPolling();
    });
  }

  void _openPostDialog() {
    showDialog<void>(
      context: context,
      builder: (_) => const PostMessageDialog(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final profile = context.watch<ProfileProvider>().profile;
    final msgProvider = context.watch<MessageProvider>();

    return Scaffold(
      appBar: _buildAppBar(context, profile),
      body: _buildBody(context, msgProvider),
      floatingActionButton: FloatingActionButton(
        onPressed: _openPostDialog,
        tooltip: 'Post a message',
        child: const Icon(Icons.edit),
      ),
    );
  }

  AppBar _buildAppBar(BuildContext context, AppProfile? profile) {
    final hasPhoto = profile != null && profile.appPhotoUrl.isNotEmpty;

    return AppBar(
      leading: hasPhoto
          ? Padding(
              padding: const EdgeInsets.all(8),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(4),
                child: Image.network(
                  profile.appPhotoUrl,
                  fit: BoxFit.contain,
                  errorBuilder: (_, __, ___) =>
                      const Icon(Icons.forum),
                ),
              ),
            )
          : const Padding(
              padding: EdgeInsets.all(8),
              child: Icon(Icons.forum),
            ),
      title: Text(profile?.appName ?? 'Messaging Board'),
      actions: [
        IconButton(
          icon: const Icon(Icons.refresh),
          tooltip: 'Refresh messages',
          onPressed: () => context.read<MessageProvider>().fetchMessages(),
        ),
      ],
    );
  }

  Widget _buildBody(BuildContext context, MessageProvider msgProvider) {
    switch (msgProvider.state) {
      case MessageState.idle:
      case MessageState.loading:
        if (msgProvider.messages.isEmpty) return const LoadingScreen();
        // Show existing messages with a progress indicator at the top
        return _buildList(msgProvider, showTopSpinner: true);

      case MessageState.error:
        if (msgProvider.messages.isEmpty) {
          return ErrorScreen(
            message: msgProvider.errorMessage,
            onRetry: () => context.read<MessageProvider>().fetchMessages(),
          );
        }
        return _buildList(msgProvider);

      case MessageState.loaded:
        if (msgProvider.messages.isEmpty) {
          return Center(
            child: Text(
              'No messages yet.\nBe the first to post!',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 16,
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.5),
              ),
            ),
          );
        }
        return _buildList(msgProvider);
    }
  }

  Widget _buildList(MessageProvider msgProvider, {bool showTopSpinner = false}) {
    return RefreshIndicator(
      onRefresh: () => context.read<MessageProvider>().fetchMessages(),
      child: CustomScrollView(
        slivers: [
          if (showTopSpinner)
            const SliverToBoxAdapter(
              child: LinearProgressIndicator(minHeight: 2),
            ),
          SliverPadding(
            padding: const EdgeInsets.only(top: 8, bottom: 80),
            sliver: SliverList.builder(
              itemCount: msgProvider.messages.length,
              itemBuilder: (context, index) =>
                  MessageCard(message: msgProvider.messages[index]),
            ),
          ),
        ],
      ),
    );
  }
}
