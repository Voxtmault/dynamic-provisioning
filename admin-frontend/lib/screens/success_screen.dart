import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../models/tenant.dart';

class SuccessScreen extends StatelessWidget {
  final RegisterTenantResponse response;

  const SuccessScreen({super.key, required this.response});

  @override
  Widget build(BuildContext context) {
    // Build the tenant's access URL from the browser's current origin.
    // e.g. admin.example.com → tenant-1.example.com
    final subdomain = response.subdomain;
    final base = Uri.base;
    // Strip the "admin" (or any first) segment from the host and prepend
    // the tenant subdomain.  e.g. admin.example.com → tenant-1.example.com
    final hostParts = base.host.split('.');
    final rootDomain = hostParts.length > 1
        ? hostParts.sublist(1).join('.')
        : base.host;
    final tenantUrl = '${base.scheme}://$subdomain.$rootDomain';

    return Scaffold(
      appBar: AppBar(title: const Text('Tenant Registered')),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 480),
          child: Card(
            margin: const EdgeInsets.all(24),
            child: Padding(
              padding: const EdgeInsets.all(32),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(
                    Icons.check_circle_outline,
                    size: 64,
                    color: Colors.green,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Tenant Provisioned Successfully',
                    style: Theme.of(context)
                        .textTheme
                        .titleLarge
                        ?.copyWith(fontWeight: FontWeight.bold),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 24),
                  _InfoRow(
                    label: 'Tenant ID',
                    value: response.tenantId.toString(),
                  ),
                  const Divider(height: 24),
                  _InfoRow(
                    label: 'Subdomain',
                    value: subdomain,
                  ),
                  const Divider(height: 24),
                  _UrlRow(label: 'Access URL', url: tenantUrl),
                  const SizedBox(height: 28),
                  Row(
                    children: [
                      Expanded(
                        child: OutlinedButton(
                          onPressed: () {
                            // Pop back to dashboard (SuccessScreen was pushed
                            // via pushReplacement, so pop goes to dashboard).
                            Navigator.of(context).pop();
                          },
                          child: const Text('Back to Dashboard'),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: FilledButton(
                          onPressed: () {
                            Navigator.of(context).pop();
                            // Dashboard's FAB will re-open register screen.
                          },
                          child: const Text('Register Another'),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final String label;
  final String value;

  const _InfoRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 100,
          child: Text(
            label,
            style: const TextStyle(
              fontWeight: FontWeight.w600,
              color: Colors.black54,
              fontSize: 13,
            ),
          ),
        ),
        Expanded(
          child: SelectableText(
            value,
            style: const TextStyle(fontSize: 14),
          ),
        ),
      ],
    );
  }
}

class _UrlRow extends StatelessWidget {
  final String label;
  final String url;

  const _UrlRow({required this.label, required this.url});

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 100,
          child: Text(
            label,
            style: const TextStyle(
              fontWeight: FontWeight.w600,
              color: Colors.black54,
              fontSize: 13,
            ),
          ),
        ),
        Expanded(
          child: Row(
            children: [
              Expanded(
                child: SelectableText(
                  url,
                  style: const TextStyle(
                    fontSize: 14,
                    color: Colors.indigo,
                    decoration: TextDecoration.underline,
                  ),
                ),
              ),
              IconButton(
                icon: const Icon(Icons.copy, size: 18),
                tooltip: 'Copy URL',
                onPressed: () {
                  Clipboard.setData(ClipboardData(text: url));
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('URL copied to clipboard')),
                  );
                },
              ),
            ],
          ),
        ),
      ],
    );
  }
}
