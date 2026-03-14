package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/iface"
	"github.com/voxtmault/dynamic-provisioning/admin-backend/internal/model"
)

type authService struct {
	userRepo  iface.AdminUserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo iface.AdminUserRepository, jwtSecret string) *authService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) SeedAdmin(ctx context.Context, email, password string) error {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		log.Println("admin user already exists, skipping seed")
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) && err.Error() != fmt.Sprintf("failed to find admin user: %s", gorm.ErrRecordNotFound.Error()) {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.AdminUser{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println("admin user seeded successfully")
	return nil
}
