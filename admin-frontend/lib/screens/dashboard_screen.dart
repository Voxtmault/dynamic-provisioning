import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../providers/auth_provider.dart';
import '../providers/tenant_provider.dart';
import '../widgets/error_screen.dart';
import '../widgets/loading_screen.dart';
import '../widgets/tenant_card.dart';
import 'register_tenant_screen.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _loadTenants());
  }

  void _loadTenants() {
    final token = context.read<AuthProvider>().token;
    if (token != null) {
      context.read<TenantProvider>().loadTenants(token);
    }
  }

  Future<void> _logout() async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Logout'),
        content: const Text('Are you sure you want to log out?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('Logout'),
          ),
        ],
      ),
    );
    if (confirm == true && mounted) {
      await context.read<AuthProvider>().logout();
    }
  }

  void _openRegister() {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => const RegisterTenantScreen(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final tenantProvider = context.watch<TenantProvider>();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Admin Panel'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            tooltip: 'Refresh',
            onPressed: _loadTenants,
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            tooltip: 'Logout',
            onPressed: _logout,
          ),
        ],
      ),
      body: _buildBody(tenantProvider),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _openRegister,
        icon: const Icon(Icons.add),
        label: const Text('Register Tenant'),
      ),
    );
  }

  Widget _buildBody(TenantProvider p) {
    switch (p.state) {
      case TenantState.idle:
      case TenantState.loading:
        if (p.tenants.isEmpty) return const LoadingScreen();
        // Show cached list with top progress bar while refreshing.
        return _buildList(p, showSpinner: true);

      case TenantState.error:
        if (p.tenants.isEmpty) {
          return ErrorScreen(
            message: p.errorMessage,
            onRetry: _loadTenants,
          );
        }
        return _buildList(p);

      case TenantState.loaded:
        if (p.tenants.isEmpty) {
          return const Center(
            child: Text(
              'No tenants registered yet.\nTap + to register the first one.',
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 15, color: Colors.black45),
            ),
          );
        }
        return _buildList(p);
    }
  }

  Widget _buildList(TenantProvider p, {bool showSpinner = false}) {
    return RefreshIndicator(
      onRefresh: () async => _loadTenants(),
      child: CustomScrollView(
        slivers: [
          if (showSpinner)
            const SliverToBoxAdapter(
              child: LinearProgressIndicator(minHeight: 2),
            ),
          SliverPadding(
            padding: const EdgeInsets.only(top: 8, bottom: 88),
            sliver: SliverList(
              delegate: SliverChildBuilderDelegate(
                (_, i) => TenantCard(tenant: p.tenants[i]),
                childCount: p.tenants.length,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
