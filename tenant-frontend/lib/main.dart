import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'providers/message_provider.dart';
import 'providers/profile_provider.dart';
import 'screens/board_screen.dart';
import 'services/api_service.dart';
import 'widgets/error_screen.dart';
import 'widgets/loading_screen.dart';

void main() {
  runApp(const TenantApp());
}

class TenantApp extends StatelessWidget {
  const TenantApp({super.key});

  @override
  Widget build(BuildContext context) {
    final apiService = ApiService();

    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => ProfileProvider(apiService)..loadProfile(),
        ),
        ChangeNotifierProvider(
          create: (_) => MessageProvider(apiService),
        ),
      ],
      child: const _AppRoot(),
    );
  }
}

class _AppRoot extends StatelessWidget {
  const _AppRoot();

  @override
  Widget build(BuildContext context) {
    final profileProvider = context.watch<ProfileProvider>();

    return MaterialApp(
      title: 'Messaging Board',
      debugShowCheckedModeBanner: false,
      theme: profileProvider.themeData,
      home: switch (profileProvider.state) {
        ProfileState.idle || ProfileState.loading => const LoadingScreen(),
        ProfileState.error => ErrorScreen(
            message: profileProvider.errorMessage,
            onRetry: () => context.read<ProfileProvider>().loadProfile(),
          ),
        ProfileState.loaded => const BoardScreen(),
      },
    );
  }
}
