package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/config"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/docker"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/secret"
)

const networkName = "saas_network"

type provisioningService struct {
	baoClient    *secret.Client
	dockerClient *docker.Client
	adminDB      *gorm.DB
	cfg          *config.Config
}

func NewProvisioningService(
	baoClient *secret.Client,
	dockerClient *docker.Client,
	adminDB *gorm.DB,
	cfg *config.Config,
) *provisioningService {
	return &provisioningService{
		baoClient:    baoClient,
		dockerClient: dockerClient,
		adminDB:      adminDB,
		cfg:          cfg,
	}
}

func (s *provisioningService) ProvisionTenant(ctx context.Context, tenant *model.Tenant) (string, string, error) {
	tenantPrefix := fmt.Sprintf("tenant_%d", tenant.ID)

	// 1. Create OpenBao policy
	policyName := fmt.Sprintf("%s-policy", tenantPrefix)
	policyRules := fmt.Sprintf(`
path "secret/data/%s/dp/%s" {
  capabilities = ["read"]
}
path "secret/metadata/%s/dp/%s" {
  capabilities = ["read"]
}
`, s.cfg.Env, tenantPrefix, s.cfg.Env, tenantPrefix)

	log.Printf("creating openbao policy: %s", policyName)
	if err := s.baoClient.CreatePolicy(policyName, policyRules); err != nil {
		return "", "", fmt.Errorf("failed to create policy: %w", err)
	}

	// 2. Create AppRole
	roleName := tenantPrefix
	log.Printf("creating approle: %s", roleName)
	if err := s.baoClient.CreateAppRole(roleName, []string{policyName}); err != nil {
		return "", "", fmt.Errorf("failed to create approle: %w", err)
	}

	// 3. Get role_id and generate secret_id
	roleID, err := s.baoClient.GetRoleID(roleName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get role-id: %w", err)
	}

	secretID, err := s.baoClient.GenerateSecretID(roleName)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate secret-id: %w", err)
	}

	// 4. Create PostgreSQL database and user
	dbName := tenantPrefix
	dbUser := fmt.Sprintf("%s_user", tenantPrefix)
	dbPassword := generateRandomPassword(32)

	log.Printf("creating database: %s", dbName)

	createDBSQL := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if err := s.adminDB.Exec(createDBSQL).Error; err != nil {
		return "", "", fmt.Errorf("failed to create database: %w", err)
	}

	createUserSQL := fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s'", dbUser, dbPassword)
	if err := s.adminDB.Exec(createUserSQL).Error; err != nil {
		return "", "", fmt.Errorf("failed to create user: %w", err)
	}

	grantDBSQL := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, dbUser)
	if err := s.adminDB.Exec(grantDBSQL).Error; err != nil {
		return "", "", fmt.Errorf("failed to grant db privileges: %w", err)
	}

	// Connect to the new database to grant schema privileges
	tenantDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.cfg.DBHost, s.cfg.DBPort, s.cfg.DBUser, s.cfg.DBPassword, dbName,
	)
	tenantDB, err := gorm.Open(
		postgres.Open(tenantDSN), &gorm.Config{},
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to tenant db: %w", err)
	}

	grantSchemaSQL := fmt.Sprintf("GRANT ALL ON SCHEMA public TO %s", dbUser)
	if err := tenantDB.Exec(grantSchemaSQL).Error; err != nil {
		return "", "", fmt.Errorf("failed to grant schema privileges: %w", err)
	}

	sqlDB, err := tenantDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// 5. Write secrets to OpenBao
	secretPath := fmt.Sprintf("%s/dp/%s", s.cfg.Env, tenantPrefix)
	secretData := map[string]interface{}{
		"db_host":        s.cfg.DBHost,
		"db_port":        s.cfg.DBPort,
		"db_user":        dbUser,
		"db_password":    dbPassword,
		"db_name":        dbName,
		"s3_endpoint":    s.cfg.S3Endpoint,
		"s3_access_key":  s.cfg.S3AccessKey,
		"s3_secret_key":  s.cfg.S3SecretKey,
		"s3_bucket":      s.cfg.S3Bucket,
		"redis_host":     s.cfg.RedisHost,
		"redis_port":     s.cfg.RedisPort,
		"redis_password": s.cfg.RedisPassword,
	}

	log.Printf("writing secrets to openbao: %s", secretPath)
	if err := s.baoClient.WriteSecret(secretPath, secretData); err != nil {
		return "", "", fmt.Errorf("failed to write secrets: %w", err)
	}

	return roleID, secretID, nil
}

func (s *provisioningService) SpinUpContainers(ctx context.Context, tenant *model.Tenant, roleID, secretID string) error {
	tenantPrefix := fmt.Sprintf("tenant_%d", tenant.ID)
	tenantSubdomain := fmt.Sprintf("tenant-%d", tenant.ID)

	// Backend container
	backendName := fmt.Sprintf("%s-backend", tenantPrefix)
	backendLabels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.%s-backend.rule", tenantPrefix):                      fmt.Sprintf("Host(`%s.localhost`) && PathPrefix(`/api`)", tenantSubdomain),
		fmt.Sprintf("traefik.http.routers.%s-backend.entrypoints", tenantPrefix):               "web",
		fmt.Sprintf("traefik.http.services.%s-backend.loadbalancer.server.port", tenantPrefix): "8080",
	}

	backendEnv := []string{
		fmt.Sprintf("OPENBAO_ADDR=%s", s.cfg.OpenBaoAddr),
		fmt.Sprintf("OPENBAO_ROLE_ID=%s", roleID),
		fmt.Sprintf("OPENBAO_SECRET_ID=%s", secretID),
		fmt.Sprintf("TENANT_ID=%s", tenantPrefix),
		fmt.Sprintf("ENV=%s", s.cfg.Env),
		fmt.Sprintf("ADMIN_BACKEND_URL=%s", s.cfg.AdminBackendURL),
	}

	log.Printf("creating backend container: %s", backendName)
	_, err := s.dockerClient.CreateAndStartContainer(ctx, docker.ContainerConfig{
		Name:        backendName,
		Image:       "tenant-backend:latest",
		Env:         backendEnv,
		Labels:      backendLabels,
		NetworkName: networkName,
	})
	if err != nil {
		return fmt.Errorf("failed to create backend container: %w", err)
	}

	// Frontend container
	frontendName := fmt.Sprintf("%s-frontend", tenantPrefix)
	frontendLabels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.%s-frontend.rule", tenantPrefix):                      fmt.Sprintf("Host(`%s.localhost`)", tenantSubdomain),
		fmt.Sprintf("traefik.http.routers.%s-frontend.entrypoints", tenantPrefix):               "web",
		fmt.Sprintf("traefik.http.services.%s-frontend.loadbalancer.server.port", tenantPrefix): "80",
		fmt.Sprintf("traefik.http.routers.%s-frontend.priority", tenantPrefix):                  "1",
	}

	log.Printf("creating frontend container: %s", frontendName)
	_, err = s.dockerClient.CreateAndStartContainer(ctx, docker.ContainerConfig{
		Name:        frontendName,
		Image:       "tenant-frontend:latest",
		Env:         []string{},
		Labels:      frontendLabels,
		NetworkName: networkName,
	})
	if err != nil {
		return fmt.Errorf("failed to create frontend container: %w", err)
	}

	return nil
}

func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("failed to generate random password: %v", err))
	}
	return hex.EncodeToString(bytes)
}
