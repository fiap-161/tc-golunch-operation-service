package usecase

import (
	"github.com/fiap-161/tc-golunch-operation-service/internal/auth/gateway"
)

type GenerateTokenUseCase struct {
	tokenGateway gateway.TokenGateway
}

func NewGenerateTokenUseCase(tokenGateway gateway.TokenGateway) *GenerateTokenUseCase {
	return &GenerateTokenUseCase{
		tokenGateway: tokenGateway,
	}
}

func (uc *GenerateTokenUseCase) Execute(userID, userType string, additionalClaims map[string]any) (string, error) {
	return uc.tokenGateway.GenerateToken(userID, userType, additionalClaims)
}
