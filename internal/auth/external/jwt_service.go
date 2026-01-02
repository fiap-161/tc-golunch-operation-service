package external

import (
	"time"

	"github.com/fiap-161/tc-golunch-operation-service/internal/auth/entity"
	"github.com/fiap-161/tc-golunch-operation-service/internal/auth/gateway"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey      string
	expiryDuration time.Duration
}

func NewJWTService(secretKey string, expiryDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:      secretKey,
		expiryDuration: expiryDuration,
	}
}

var _ gateway.TokenGateway = (*JWTService)(nil)

func (s *JWTService) GenerateToken(userID, userType string, additionalClaims map[string]any) (string, error) {
	now := time.Now()

	claims := entity.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:   userID,
		UserType: userType,
		Custom:   additionalClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*entity.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*entity.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
