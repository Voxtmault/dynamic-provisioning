class TenantListItem {
  final int id;
  final String name;
  final String subdomain;
  final String status;
  final String backendContainerStatus;
  final String frontendContainerStatus;

  const TenantListItem({
    required this.id,
    required this.name,
    required this.subdomain,
    required this.status,
    required this.backendContainerStatus,
    required this.frontendContainerStatus,
  });

  factory TenantListItem.fromJson(Map<String, dynamic> json) {
    return TenantListItem(
      id: (json['id'] as num).toInt(),
      name: json['name'] as String? ?? '',
      subdomain: json['subdomain'] as String? ?? '',
      status: json['status'] as String? ?? 'unknown',
      backendContainerStatus:
          json['backend_container_status'] as String? ?? 'unknown',
      frontendContainerStatus:
          json['frontend_container_status'] as String? ?? 'unknown',
    );
  }
}

class RegisterTenantResponse {
  final int tenantId;
  final String subdomain;
  final String photoUploadUrl;

  const RegisterTenantResponse({
    required this.tenantId,
    required this.subdomain,
    required this.photoUploadUrl,
  });

  factory RegisterTenantResponse.fromJson(Map<String, dynamic> json) {
    return RegisterTenantResponse(
      tenantId: (json['tenant_id'] as num).toInt(),
      subdomain: json['subdomain'] as String? ?? '',
      photoUploadUrl: json['photo_upload_url'] as String? ?? '',
    );
  }
}
