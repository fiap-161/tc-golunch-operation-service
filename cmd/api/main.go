package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/fiap-161/tc-golunch-operation-service/database"
	_ "github.com/fiap-161/tc-golunch-operation-service/docs"

	// admincontroller "github.com/fiap-161/tc-golunch-operation-service/internal/admin/controller"
	// adminmodel "github.com/fiap-161/tc-golunch-operation-service/internal/admin/dto"
	// admindatasource "github.com/fiap-161/tc-golunch-operation-service/internal/admin/external/datasource"
	// admingateway "github.com/fiap-161/tc-golunch-operation-service/internal/admin/gateway"
	// adminhandler "github.com/fiap-161/tc-golunch-operation-service/internal/admin/handler"
	authcontroller "github.com/fiap-161/tc-golunch-operation-service/internal/auth/controller"
	"github.com/fiap-161/tc-golunch-operation-service/internal/auth/external"
	"github.com/fiap-161/tc-golunch-operation-service/internal/http/middleware"
	ordercontroller "github.com/fiap-161/tc-golunch-operation-service/internal/order/controller"
	ordermodel "github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	orderdatasource "github.com/fiap-161/tc-golunch-operation-service/internal/order/external/datasource"
	ordergateway "github.com/fiap-161/tc-golunch-operation-service/internal/order/gateway"
	orderhandler "github.com/fiap-161/tc-golunch-operation-service/internal/order/handler"
	orderusecases "github.com/fiap-161/tc-golunch-operation-service/internal/order/usecases"
	"github.com/fiap-161/tc-golunch-operation-service/internal/shared/httpclient"
)

// @title           GoLunch Operation Service API
// @version         1.0
// @description     API para gerenciamento das operações da cozinha e painel administrativo da lanchonete GoLunch
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host            localhost:8083
// @BasePath        /
func main() {
	r := gin.Default()
	loadYAML()

	db := database.NewPostgresDatabase().GetDb()

	if err := db.AutoMigrate(
		// &adminmodel.AdminDAO{}, // REMOVED: Admin now in Auth Service
		&ordermodel.OrderDAO{},
	); err != nil {
		log.Fatalf("Erro ao migrar o banco: %v", err)
	}

	// JWT service for generate and validate tokens
	jwtGateway := external.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	authController := authcontroller.New(jwtGateway)

	// Admin - REMOVED: now centralized in Auth Service (Order Service)
	// adminDatasource := admindatasource.New(db)
	// authGateway := admingateway.NewAuthGateway(authController)
	// adminController := admincontroller.Build(adminDatasource, authGateway)
	// adminHandler := adminhandler.New(adminController)

	// Order Data Source and Gateway
	orderDataSource := orderdatasource.New(db)
	orderGateway := ordergateway.Build(orderDataSource)

	// HTTP Clients para outros serviços
	productClient := httpclient.NewProductClient("http://localhost:8081")
	productOrderClient := httpclient.NewProductOrderClient("http://localhost:8081")
	paymentClient := httpclient.NewPaymentClient("http://localhost:8082")

	// Order Use Case
	orderUseCase := orderusecases.Build(orderGateway, productClient, productOrderClient, paymentClient)

	// Order Controller and Handler
	orderController := ordercontroller.Build(orderUseCase)
	orderHandler := orderhandler.New(orderController)

	// Default Routes
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	// Public Routes - REMOVED: auth now handled by Auth Service (Order Service port 8081)
	// r.POST("/admin/register", adminHandler.Register)
	// r.POST("/admin/login", adminHandler.Login)
	// r.GET("/admin/validate", adminHandler.ValidateToken)

	// Auth Service endpoints documentation
	r.GET("/auth-info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":          "Admin authentication is now centralized",
			"auth_service_url": "http://localhost:8081",
			"admin_register":   "POST http://localhost:8081/admin/register",
			"admin_login":      "POST http://localhost:8081/admin/login",
			"admin_validate":   "POST http://localhost:8081/admin/validate",
		})
	})

	// Authenticated Group
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware(authController))

	// Admin Routes
	adminRoutes := authenticated.Group("/admin")
	adminRoutes.Use(middleware.AdminOnly())

	// Order Management Routes
	adminRoutes.GET("/orders", orderHandler.GetAll)
	adminRoutes.PUT("/orders/:id", orderHandler.Update)
	adminRoutes.GET("/orders/panel", orderHandler.GetPanel)

	r.Run(":8083")
}

func loadYAML() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/environment")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading yaml config: %v", err)
	}
}

// Ping godoc
// @Summary      Answers with "pong"
// @Description  Health Check
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200 {object}  PongResponse
// @Router       /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type PongResponse struct {
	Message string `json:"message"`
}
