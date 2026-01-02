package usecases

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/entity"
	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/gateway"
	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/utils"
	apperror "github.com/fiap-161/tc-golunch-operation-service/internal/shared/errors"
)

type UseCases struct {
	AdminGateway gateway.Gateway
}

func Build(productGateway gateway.Gateway) *UseCases {
	return &UseCases{AdminGateway: productGateway}
}

func (u *UseCases) Create(ctx context.Context, admin entity.Admin) error {

	saved, _ := u.FindByEmail(ctx, admin.Email)
	if saved.Email != "" {
		return &apperror.ValidationError{Msg: "User already registered"}
	}

	hash, err := utils.HashPassword(admin.Password)

	if err != nil {
		return err
	}

	adminHashed := admin.Build(hash)

	err = u.AdminGateway.Create(ctx, adminHashed)

	if err != nil {
		return err
	}

	return nil
}

func (u *UseCases) FindByEmail(ctx context.Context, email string) (entity.Admin, error) {

	admin, err := u.AdminGateway.FindByEmail(ctx, email)
	if err != nil {
		return entity.Admin{}, err
	}

	return admin, nil
}

func (u *UseCases) Login(ctx context.Context, admin entity.Admin) (string, bool, error) {

	saved, err := u.FindByEmail(ctx, admin.Email)
	if err != nil {
		return "", true, &apperror.UnauthorizedError{Msg: "Invalid email or password"}
	}

	if !utils.CheckPasswordHash(admin.Password, saved.Password) {
		return "", true, &apperror.UnauthorizedError{Msg: "Invalid email or password"}
	}

	return saved.Id, true, nil
}
