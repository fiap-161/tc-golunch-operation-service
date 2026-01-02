package gateway

import (
	authcontroller "github.com/fiap-161/tc-golunch-operation-service/internal/auth/controller"
)

// AuthGatewayImpl implementa AuthGateway usando o auth controller
type AuthGatewayImpl struct {
	authController *authcontroller.Controller
}

func NewAuthGateway(authController *authcontroller.Controller) AuthGateway {
	return &AuthGatewayImpl{
		authController: authController,
	}
}

func (a *AuthGatewayImpl) GenerateToken(userID string, userType string, additionalClaims map[string]any) (string, error) {
	return a.authController.GenerateToken(userID, userType, additionalClaims)
}

func (a *AuthGatewayImpl) ValidateToken(token string) (map[string]interface{}, error) {
	claims, err := a.authController.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Convert claims to map[string]interface{}
	result := make(map[string]interface{})
	result["user_id"] = claims.UserID
	result["role"] = claims.UserType

	return result, nil
}
