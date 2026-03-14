package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/controller"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/repo"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/router"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/service"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/config"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/store/db/postgre"
	redisStore "github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/store/db/redis"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/store/object"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/store/secret"
)

func main() {
	// 1. Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Println("configuration loaded")

	// 2. Connect to OpenBao and retrieve tenant secrets
	baoClient, err := secret.NewClient(cfg.OpenBaoAddr)
	if err != nil {
		log.Fatalf("failed to create openbao client: %v", err)
	}

	if err := baoClient.LoginAppRole(cfg.OpenBaoRoleID, cfg.OpenBaoSecretID); err != nil {
		log.Fatalf("failed to login to openbao: %v", err)
	}
	log.Println("authenticated to openbao")

	secrets, err := baoClient.ReadTenantSecrets(cfg.Env, cfg.TenantID)
	if err != nil {
		log.Fatalf("failed to read tenant secrets: %v", err)
	}
	log.Println("tenant secrets retrieved")

	// 3. Connect to PostgreSQL and run migrations
	db, err := postgre.NewConnection(
		secrets.DBHost, secrets.DBPort,
		secrets.DBUser, secrets.DBPassword, secrets.DBName,
	)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	log.Println("connected to postgres")

	if err := postgre.RunMigrations(db, &model.Message{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("database migrations completed")

	// 4. Connect to Redis
	redisClient, err := redisStore.NewClient(
		secrets.RedisHost, secrets.RedisPort, secrets.RedisPassword,
	)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("connected to redis")

	// 5. Initialize S3/Garage client
	s3Client, err := object.NewS3Client(
		secrets.S3Endpoint, secrets.S3AccessKey, secrets.S3SecretKey,
		secrets.S3Bucket, false,
	)
	if err != nil {
		log.Fatalf("failed to create s3 client: %v", err)
	}
	log.Println("s3 client initialized")

	// 6. Initialize profile service and fetch initial profile
	profileSvc := service.NewProfileService(redisClient, s3Client, cfg.AdminBackendURL, cfg.TenantID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := profileSvc.InitProfile(ctx); err != nil {
		log.Printf("warning: failed to initialize profile at startup: %v (will retry on first request)", err)
	}
	cancel()

	// 7. Wire dependencies
	messageRepo := repo.NewMessageRepository(db)
	messageSvc := service.NewMessageService(messageRepo)

	msgCtrl := controller.NewMessageController(messageSvc)
	profileCtrl := controller.NewProfileController(profileSvc)

	// 8. Setup Echo and routes
	e := echo.New()
	e.HideBanner = true
	router.Setup(e, msgCtrl, profileCtrl)

	// 9. Start server with graceful shutdown
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

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		log.Printf("warning: failed to close redis connection: %v", err)
	}

	// Close Postgres connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("server exited")
}
