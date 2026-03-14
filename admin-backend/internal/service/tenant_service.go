package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/docker"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/object"
)

const presignedPutExpiry = 30 * time.Minute

type tenantService struct {
	tenantRepo      iface.TenantRepository
	provisioningSvc iface.ProvisioningService
	s3Client        *object.S3Client
	dockerClient    *docker.Client
}

func NewTenantService(
	tenantRepo iface.TenantRepository,
	provisioningSvc iface.ProvisioningService,
	s3Client *object.S3Client,
	dockerClient *docker.Client,
) *tenantService {
	return &tenantService{
		tenantRepo:      tenantRepo,
		provisioningSvc: provisioningSvc,
		s3Client:        s3Client,
		dockerClient:    dockerClient,
	}
}

func (s *tenantService) RegisterTenant(ctx context.Context, req model.RegisterTenantRequest) (*model.RegisterTenantResponse, error) {
	// 1. Create tenant record
	tenant := &model.Tenant{
		Name:         req.Name,
		ColorPalette: req.ColorPalette,
		Status:       model.TenantStatusProvisioning,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant record: %w", err)
	}

	// Set subdomain and photo key after we have the ID
	tenant.Subdomain = fmt.Sprintf("tenant-%d", tenant.ID)
	tenant.AppPhotoKey = fmt.Sprintf("tenants/tenant_%d/app_photo", tenant.ID)

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant with subdomain: %w", err)
	}

	// 2. Provision infrastructure
	log.Printf("provisioning tenant %d: %s", tenant.ID, tenant.Name)
	roleID, secretID, err := s.provisioningSvc.ProvisionTenant(ctx, tenant)
	if err != nil {
		tenant.Status = model.TenantStatusError
		_ = s.tenantRepo.Update(ctx, tenant)
		return nil, fmt.Errorf("failed to provision tenant: %w", err)
	}

	// 3. Spin up containers
	log.Printf("spinning up containers for tenant %d", tenant.ID)
	if err := s.provisioningSvc.SpinUpContainers(ctx, tenant, roleID, secretID); err != nil {
		tenant.Status = model.TenantStatusError
		_ = s.tenantRepo.Update(ctx, tenant)
		return nil, fmt.Errorf("failed to spin up containers: %w", err)
	}

	// 4. Mark tenant as active
	tenant.Status = model.TenantStatusActive
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant status: %w", err)
	}

	// 5. Generate presigned PUT URL for photo upload
	uploadURL, err := s.s3Client.GeneratePresignedPutURL(tenant.AppPhotoKey, presignedPutExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload url: %w", err)
	}

	log.Printf("tenant %d provisioned successfully", tenant.ID)

	return &model.RegisterTenantResponse{
		TenantID:       tenant.ID,
		Subdomain:      tenant.Subdomain,
		PhotoUploadURL: uploadURL,
	}, nil
}

func (s *tenantService) ListTenants(ctx context.Context) ([]model.TenantListItem, error) {
	tenants, err := s.tenantRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	items := make([]model.TenantListItem, 0, len(tenants))
	for _, t := range tenants {
		tenantPrefix := fmt.Sprintf("tenant_%d", t.ID)

		backendStatus, _ := s.dockerClient.GetContainerStatus(ctx, tenantPrefix+"-backend")
		frontendStatus, _ := s.dockerClient.GetContainerStatus(ctx, tenantPrefix+"-frontend")

		item := model.TenantListItem{
			ID:        t.ID,
			Name:      t.Name,
			Subdomain: t.Subdomain,
			Status:    t.Status,
		}

		if backendStatus != nil {
			item.BackendContainerStatus = backendStatus.State
		} else {
			item.BackendContainerStatus = "unknown"
		}

		if frontendStatus != nil {
			item.FrontendContainerStatus = frontendStatus.State
		} else {
			item.FrontendContainerStatus = "unknown"
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *tenantService) GetTenantProfile(ctx context.Context, tenantID uint) (*model.TenantProfileResponse, error) {
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant: %w", err)
	}

	return &model.TenantProfileResponse{
		TenantID:     fmt.Sprintf("tenant_%d", tenant.ID),
		AppName:      tenant.Name,
		AppPhotoKey:  tenant.AppPhotoKey,
		ColorPalette: tenant.ColorPalette,
	}, nil
}
