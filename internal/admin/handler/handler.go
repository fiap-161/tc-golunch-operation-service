package handler

import (
	"context"
	"net/http"

	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/controller"
	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/dto"
	apperror "github.com/fiap-161/tc-golunch-operation-service/internal/shared/errors"
	"github.com/fiap-161/tc-golunch-operation-service/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminController *controller.Controller
}

func New(adminController *controller.Controller) *Handler {
	return &Handler{
		adminController: adminController,
	}
}

// Register godoc
// @Summary      Register Admin
// @Description  Register a new admin user
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AdminRequestDTO  true  "Admin registration details"
// @Success      201      {object}  map[string]string     "Success message"
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /admin/register [post]
func (h *Handler) Register(c *gin.Context) {
	ctx := context.Background()

	var adminRequest dto.AdminRequestDTO
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	err := h.adminController.Register(ctx, adminRequest)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// Login godoc
// @Summary      Admin Login
// @Description  Authenticates an admin user and returns a JWT token
// @Tags         Admin Domain
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AdminRequestDTO  true  "Admin login credentials"
// @Success      200      {object}  TokenDTO
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      401      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /admin/login [post]
func (h *Handler) Login(c *gin.Context) {
	ctx := context.Background()

	var adminRequest dto.AdminRequestDTO
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	token, err := h.adminController.Login(ctx, adminRequest)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &TokenDTO{
		TokenString: token,
	})
}

type TokenDTO struct {
	TokenString string `json:"token"`
}

// ValidateToken godoc
// @Summary      Validate Admin Token
// @Description  Validates an admin JWT token for inter-service communication
// @Tags         Admin Domain
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]interface{}  "Token validation result"
// @Failure      401      {object}  errors.ErrorDTO
// @Router       /admin/validate [get]
func (h *Handler) ValidateToken(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required", "valid": false})
		return
	}

	// Parse Bearer token
	tokenParts := c.Request.Header.Get("Authorization")
	if len(tokenParts) < 7 || tokenParts[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format", "valid": false})
		return
	}

	token := tokenParts[7:]

	// Validate token using admin controller
	ctx := context.Background()
	isValid, adminData := h.adminController.ValidateToken(ctx, token)

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"admin": adminData,
	})
}
