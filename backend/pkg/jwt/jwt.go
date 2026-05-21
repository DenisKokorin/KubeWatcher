package jwt

import (
	"errors"
	"fmt"
	"k8s-mon/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrClaimsFailed         = errors.New("claims failed")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type JWTService struct {
	Secret    string
	AccessTTL time.Duration
}

type TokenService interface {
	GenerateAccessToken(user models.Employee) (string, error)
	RegenerateToken(oldToken string) (string, error)
	ParseToken(token string) (*UserClaims, error)
}

func NewUserJWTpkg(Secret string, AccessTTL time.Duration) JWTService {
	return JWTService{Secret: Secret, AccessTTL: AccessTTL}
}

func (j JWTService) GenerateAccessToken(user models.Employee) (string, error) {
	const op = "jwt.GenerateAccessToken"

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTTL)),
		},
		UUID:     user.UUID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Team:     user.Team,
		Position: user.Position,
		Role:     user.Role,
	}

	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	AccessTokenString, err := AccessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessTokenString, nil
}

func (u JWTService) RegenerateToken(oldToken string) (string, error) {
	const op = "jwt.RegenerateToken"

	var claims UserClaims
	token, err := jwt.ParseWithClaims(oldToken, &claims, func(t *jwt.Token) (interface{}, error) { return []byte(u.Secret), nil }, jwt.WithoutClaimsValidation())
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	parsedUUID, err := uuid.Parse(claims.UUID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return "", ErrClaimsFailed
	}

	user := models.Employee{
		UUID:     parsedUUID,
		Name:     claims.Name,
		Team:     claims.Team,
		Position: claims.Position,
		Email:    claims.Email,
		Role:     claims.Role,
	}
	return u.GenerateAccessToken(user)
}

func (j JWTService) ParseToken(token string) (*UserClaims, error) {
	const op = "jwt.ParseToken"
	var claims UserClaims

	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !tkn.Valid {
		return nil, ErrClaimsFailed
	}

	return &claims, nil
}
