package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/voxtmault/dynamic-provisioning/tenant-backend/internal/model"
	"github.com/voxtmault/dynamic-provisioning/tenant-backend/pkg/store/object"
)

const (
	profileCacheKey = "app_profile"
	profileCacheTTL = 30 * time.Minute
	presignedExpiry = 30 * time.Minute
)

type profileService struct {
	redis           *redis.Client
	s3              *object.S3Client
	adminBackendURL string
	tenantID        string
}

func NewProfileService(
	redisClient *redis.Client,
	s3Client *object.S3Client,
	adminBackendURL string,
	tenantID string,
) *profileService {
	return &profileService{
		redis:           redisClient,
		s3:              s3Client,
		adminBackendURL: adminBackendURL,
		tenantID:        tenantID,
	}
}

// InitProfile fetches the profile from the admin backend and caches it.
// Should be called during startup.
func (s *profileService) InitProfile(ctx context.Context) error {
	_, err := s.fetchAndCacheProfile(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize profile: %w", err)
	}
	log.Println("tenant profile cached successfully")
	return nil
}

func (s *profileService) GetProfile(ctx context.Context) (*model.AppProfile, error) {
	// Try cache first
	cached, err := s.redis.Get(ctx, profileCacheKey).Result()
	if err == nil {
		var profile model.AppProfile
		if err := json.Unmarshal([]byte(cached), &profile); err == nil {
			return &profile, nil
		}
	}

	// Cache miss — re-fetch and cache
	return s.fetchAndCacheProfile(ctx)
}

func (s *profileService) fetchAndCacheProfile(ctx context.Context) (*model.AppProfile, error) {
	profile, err := s.fetchFromAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// Generate pre-signed URL for the app photo if a key is present
	if profile.AppPhotoKey != "" {
		presignedURL, err := s.s3.GeneratePresignedURL(profile.AppPhotoKey, presignedExpiry)
		if err != nil {
			log.Printf("warning: failed to generate presigned url for app photo: %v", err)
		} else {
			profile.AppPhotoURL = presignedURL
		}
	}

	// Cache in Redis
	data, err := json.Marshal(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := s.redis.Set(ctx, profileCacheKey, data, profileCacheTTL).Err(); err != nil {
		log.Printf("warning: failed to cache profile in redis: %v", err)
	}

	return profile, nil
}

func (s *profileService) fetchFromAdmin(ctx context.Context) (*model.AppProfile, error) {
	url := fmt.Sprintf("%s/api/tenant/%s/profile", s.adminBackendURL, strings.TrimPrefix(s.tenantID, "tenant_"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile from admin backend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("admin backend returned %d: %s", resp.StatusCode, string(body))
	}

	// We expect the admin backend to return the profile inside an APIResponse wrapper
	var apiResp struct {
		Data model.AppProfile `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode admin response: %w", err)
	}

	return &apiResp.Data, nil
}
