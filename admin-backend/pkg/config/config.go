package config

import (
	"fmt"
	"os"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// OpenBao
	OpenBaoAddr      string
	OpenBaoRootToken string

	// S3 / Garage
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Auth
	JWTSecret     string
	AdminEmail    string
	AdminPassword string

	// Application
	Env             string
	AdminBackendURL string
}

func Load() (*Config, error) {
	cfg := &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		OpenBaoAddr:      os.Getenv("OPENBAO_ADDR"),
		OpenBaoRootToken: os.Getenv("OPENBAO_ROOT_TOKEN"),
		S3Endpoint:       os.Getenv("S3_ENDPOINT"),
		S3AccessKey:      os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:      os.Getenv("S3_SECRET_KEY"),
		S3Bucket:         os.Getenv("S3_BUCKET"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		AdminEmail:       os.Getenv("ADMIN_EMAIL"),
		AdminPassword:    os.Getenv("ADMIN_PASSWORD"),
		Env:              os.Getenv("ENV"),
		AdminBackendURL:  os.Getenv("ADMIN_BACKEND_URL"),
	}

	required := map[string]string{
		"DB_HOST":            cfg.DBHost,
		"DB_PORT":            cfg.DBPort,
		"DB_USER":            cfg.DBUser,
		"DB_PASSWORD":        cfg.DBPassword,
		"DB_NAME":            cfg.DBName,
		"OPENBAO_ADDR":       cfg.OpenBaoAddr,
		"OPENBAO_ROOT_TOKEN": cfg.OpenBaoRootToken,
		"S3_ENDPOINT":        cfg.S3Endpoint,
		"S3_ACCESS_KEY":      cfg.S3AccessKey,
		"S3_SECRET_KEY":      cfg.S3SecretKey,
		"S3_BUCKET":          cfg.S3Bucket,
		"REDIS_HOST":         cfg.RedisHost,
		"REDIS_PORT":         cfg.RedisPort,
		"REDIS_PASSWORD":     cfg.RedisPassword,
		"JWT_SECRET":         cfg.JWTSecret,
		"ADMIN_EMAIL":        cfg.AdminEmail,
		"ADMIN_PASSWORD":     cfg.AdminPassword,
		"ENV":                cfg.Env,
		"ADMIN_BACKEND_URL":  cfg.AdminBackendURL,
	}

	for name, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return cfg, nil
}
