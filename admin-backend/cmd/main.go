package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/controller"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/repo"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/router"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/service"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/config"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/db/postgre"
	redisStore "github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/db/redis"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/docker"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/object"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/pkg/store/secret"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Println("configuration loaded")

	// 2. Connect to PostgreSQL (admin database)
	db, err := postgre.NewConnection(
		cfg.DBHost, cfg.DBPort,
		cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	log.Println("connected to postgres")

	// 3. Run migrations
	if err := postgre.RunMigrations(db, &model.AdminUser{}, &model.Tenant{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("database migrations completed")

	// 4. Connect to Redis
	redisClient, err := redisStore.NewClient(
		cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword,
	)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("connected to redis")

	// 5. Initialize S3/Garage client
	s3Client, err := object.NewS3Client(
		cfg.S3Endpoint, cfg.S3PublicEndpoint, cfg.S3AccessKey, cfg.S3SecretKey,
		cfg.S3Bucket, false,
	)
	if err != nil {
		log.Fatalf("failed to create s3 client: %v", err)
	}
	log.Println("s3 client initialized")

	// 6. Initialize OpenBao client (root token)
	baoClient, err := secret.NewClient(cfg.OpenBaoAddr, cfg.OpenBaoRootToken)
	if err != nil {
		log.Fatalf("failed to create openbao client: %v", err)
	}
	log.Println("openbao client initialized")

	// 7. Initialize Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("failed to create docker client: %v", err)
	}
	log.Println("docker client initialized")

	// 8. Wire repositories
	userRepo := repo.NewAdminUserRepository(db)
	tenantRepo := repo.NewTenantRepository(db)

	// 9. Wire services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	provisioningSvc := service.NewProvisioningService(baoClient, dockerClient, db, cfg)
	tenantSvc := service.NewTenantService(tenantRepo, provisioningSvc, s3Client, dockerClient)

	// 10. Seed admin user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := authSvc.SeedAdmin(ctx, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}
	cancel()

	// 11. Wire controllers
	authCtrl := controller.NewAuthController(authSvc)
	tenantCtrl := controller.NewTenantController(tenantSvc)

	// 12. Setup Echo and routes
	e := echo.New()
	e.HideBanner = true
	router.Setup(e, authCtrl, tenantCtrl, cfg.JWTSecret)

	// 13. Start server with graceful shutdown
	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	// Close connections
	if err := redisClient.Close(); err != nil {
		log.Printf("warning: failed to close redis connection: %v", err)
	}

	if err := dockerClient.Close(); err != nil {
		log.Printf("warning: failed to close docker client: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("server exited")
}
