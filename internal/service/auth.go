package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/pkg/errors"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash, name, role string, apiKey *string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, id uuid.UUID, name, role string) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type AuthService struct {
	repo      UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewAuthService(repo UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewAppError("UNAUTHORIZED", "invalid credentials", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.NewAppError("UNAUTHORIZED", "invalid credentials", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.NewAppError("INTERNAL", "failed to generate token", err)
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewAppError("INTERNAL", "failed to hash password", err)
	}

	role := req.Role
	if role == "" {
		role = model.RoleMember
	}

	apiKey := generateAPIKey()

	user, err := s.repo.Create(ctx, req.Email, string(hash), req.Name, string(role), &apiKey)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.ErrBadRequest
	}
	return s.repo.GetByID(ctx, id)
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (string, string, error) {
	// Try JWT first
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err == nil && token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return "", "", errors.ErrUnauthorized
		}
		userID, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)
		return userID, role, nil
	}

	// Try API key
	user, err := s.repo.GetByAPIKey(ctx, tokenStr)
	if err != nil {
		return "", "", errors.ErrUnauthorized
	}

	return user.ID.String(), string(user.Role), nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.jwtExpiry).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func generateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
