package handler

import (
	"context"
	"net/http"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/controller"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity/enum"
	apperror "github.com/fiap-161/tc-golunch-operation-service/internal/shared/errors"
	"github.com/fiap-161/tc-golunch-operation-service/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	controller *controller.Controller
}

func New(controller *controller.Controller) *Handler {
	return &Handler{controller: controller}
}

// Create Order godoc
// @Summary      Create Order
// @Description  Create a new order
// @Tags         Order Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateOrderDTO true "Order to create. Note that the customer_id is automatically set from the authenticated user."
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /order/ [post]
func (h *Handler) Create(c *gin.Context) {
	var orderDTO dto.CreateOrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "invalid request body",
			MessageError: err.Error(),
		})
		return
	}
	if err := orderDTO.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "validation failed",
			MessageError: err.Error(),
		})
		return
	}
	customerIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apperror.ErrorDTO{
			Message:      "unauthorized",
			MessageError: "user id not found in context",
		})
		return
	}
	customerID := customerIDRaw.(string)
	orderDTO.CustomerID = customerID

	qrCode, err := h.controller.Create(context.Background(), orderDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"qr_code": qrCode,
		"message": "Order created successfully",
	})
}

// Update Order godoc
// @Summary      Update Order
// @Description  Update an existing order status
// @Tags         Order Domain
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Param        request body dto.UpdateOrderDTO true "Order status update"
// @Success      204  "No Content"
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Failure      404  {object}  errors.ErrorDTO
// @Router       /order/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var orderUpdate dto.UpdateOrderDTO
	if err := c.ShouldBindJSON(&orderUpdate); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}
	orderDAO, err := h.controller.FindByID(context.Background(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	orderDAO.Status = enum.OrderStatus(orderUpdate.Status)
	_, err = h.controller.Update(context.Background(), orderDAO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetAll godoc
// @Summary      Get all orders
// @Description  Retrieve a list of all orders, optionally filtered by ID
// @Tags         Order Domain
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   query     string  false  "Optional order ID filter"
// @Success      200  {object}  dto.OrderResponseListDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Failure      500  {object}  errors.ErrorDTO
// @Router       /order/ [get]
func (h *Handler) GetAll(c *gin.Context) {
	id := c.Query("id")
	orders, err := h.controller.GetAll(context.Background(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.OrderResponseListDTO{
		Orders: orders,
	})
}

// GetPanel Get Order Panel godoc
// @Summary      Get Order Panel
// @Description  Get the order panel with all orders that are in the panel status
// @Tags         Order Domain
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.OrderPanelDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /order/panel [get]
func (h *Handler) GetPanel(c *gin.Context) {
	orders, err := h.controller.GetPanel(context.Background())
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	panel := dto.OrderPanelDTO{Orders: []dto.OrderPanelItemDTO{}}
	for _, order := range orders {
		panel.Orders = append(panel.Orders, dto.OrderPanelItemDTO{
			OrderNumber:   order.Entity.ID[len(order.Entity.ID)-4:],
			Status:        string(order.Status),
			PreparingTime: order.PreparingTime,
			CreatedAt:     order.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, panel)
}
