import 'package:flutter/material.dart';

import '../models/tenant.dart';

class TenantCard extends StatelessWidget {
  final TenantListItem tenant;

  const TenantCard({super.key, required this.tenant});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    tenant.name,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                _StatusChip(status: tenant.status),
              ],
            ),
            const SizedBox(height: 6),
            Text(
              tenant.subdomain,
              style: TextStyle(fontSize: 13, color: Colors.grey.shade600),
            ),
            const SizedBox(height: 10),
            Row(
              children: [
                _ContainerBadge(
                  label: 'Backend',
                  status: tenant.backendContainerStatus,
                ),
                const SizedBox(width: 8),
                _ContainerBadge(
                  label: 'Frontend',
                  status: tenant.frontendContainerStatus,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _StatusChip extends StatelessWidget {
  final String status;

  const _StatusChip({required this.status});

  Color get _color => switch (status.toLowerCase()) {
        'active' => Colors.green,
        'provisioning' => Colors.orange,
        'stopped' => Colors.grey,
        'error' => Colors.red,
        _ => Colors.blueGrey,
      };

  @override
  Widget build(BuildContext context) {
    return Chip(
      label: Text(
        status,
        style: const TextStyle(color: Colors.white, fontSize: 12),
      ),
      backgroundColor: _color,
      padding: EdgeInsets.zero,
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
    );
  }
}

class _ContainerBadge extends StatelessWidget {
  final String label;
  final String status;

  const _ContainerBadge({required this.label, required this.status});

  Color get _color => switch (status.toLowerCase()) {
        'running' => Colors.green.shade700,
        'exited' || 'stopped' => Colors.grey,
        'unknown' => Colors.blueGrey,
        _ => Colors.orange,
      };

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: _color.withValues(alpha: 0.12),
        border: Border.all(color: _color.withValues(alpha: 0.5)),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        '$label: $status',
        style: TextStyle(fontSize: 11, color: _color),
      ),
    );
  }
}
