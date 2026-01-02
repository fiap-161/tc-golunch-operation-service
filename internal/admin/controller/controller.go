package controller

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/dto"
	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/external/datasource"
	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/gateway"
	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/usecases"
)

type Controller struct {
	AdminDatasource datasource.DataSource
	AuthGateway     gateway.AuthGateway
}

func Build(productDataSource datasource.DataSource, authGateway gateway.AuthGateway) *Controller {
	return &Controller{
		AdminDatasource: productDataSource,
		AuthGateway:     authGateway,
	}
}

func (c *Controller) Register(ctx context.Context, adminRequest dto.AdminRequestDTO) error {
	adminGateway := gateway.Build(c.AdminDatasource)
	useCase := usecases.Build(*adminGateway)
	admin := dto.FromAdminRequestDTO(adminRequest)
	err := useCase.Create(ctx, admin)

	if err != nil {
		return err
	}

	return nil

}

func (c *Controller) Login(ctx context.Context, adminRequest dto.AdminRequestDTO) (string, error) {
	adminGateway := gateway.Build(c.AdminDatasource)
	useCase := usecases.Build(*adminGateway)
	admin := dto.FromAdminRequestDTO(adminRequest)
	adminId, _, err := useCase.Login(ctx, admin)

	if err != nil {
		return "", err
	}

	token, err2 := c.AuthGateway.GenerateToken(adminId, "admin", nil)

	if err2 != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken validates a JWT token and returns admin information
func (c *Controller) ValidateToken(ctx context.Context, token string) (bool, map[string]interface{}) {
	// Use AuthGateway to validate token (which is the auth controller)
	claims, err := c.AuthGateway.ValidateToken(token)
	if err != nil {
		return false, nil
	}

	// Check if the token is for admin role
	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		return false, nil
	}

	// Extract admin info from claims
	adminData := map[string]interface{}{
		"id":   claims["user_id"],
		"role": claims["role"],
	}

	return true, adminData
}
