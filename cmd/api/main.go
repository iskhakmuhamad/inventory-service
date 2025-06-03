package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/iskhakmuhamad/inventory-service/configs"
	"github.com/iskhakmuhamad/inventory-service/internal/app/controllers"
	"github.com/iskhakmuhamad/inventory-service/internal/app/middleware"
	"github.com/iskhakmuhamad/inventory-service/internal/app/repositories"
	"github.com/iskhakmuhamad/inventory-service/internal/app/services"
	"github.com/iskhakmuhamad/inventory-service/pkg/database"
	"github.com/iskhakmuhamad/inventory-service/pkg/validators"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

var validate *validator.Validate

func init() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	validate = validator.New()
	validators.RegisterCustomValidators(validate)
}

func main() {
	// Load configuration
	cfg := configs.Load()

	// Initialize database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	// if err := database.RunMigrations(db); err != nil {
	// 	log.Fatal("Failed to run migrations:", err)
	// }

	// Initialize repositories
	adminRepo := repositories.NewAdminRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize services
	adminService := services.NewAdminService(adminRepo, cfg.JWT.Secret)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo)
	transactionService := services.NewTransactionService(transactionRepo, db)

	// Initialize controllers
	adminController := controllers.NewAdminController(adminService, validate)
	categoryController := controllers.NewCategoryController(categoryService, validate)
	productController := controllers.NewProductController(productService, validate)
	transactionController := controllers.NewTransactionController(transactionService, validate)

	// Initialize router
	r := gin.Default()

	// Setup routes
	setupRoutes(r, adminController, categoryController, productController, transactionController, cfg.JWT.Secret)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}

func setupRoutes(r *gin.Engine, adminCtrl *controllers.AdminController, categoryCtrl *controllers.CategoryController,
	productCtrl *controllers.ProductController, transactionCtrl *controllers.TransactionController, jwtSecret string) {

	api := r.Group("/api/v1")

	// Public routes
	api.POST("/admin/register", adminCtrl.CreateAdmin)
	api.POST("/admin/login", adminCtrl.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Admin routes
		protected.PUT("/profile", adminCtrl.UpdateProfile)
		protected.GET("/profile", adminCtrl.GetProfile)
		protected.GET("/admin", adminCtrl.GetAllAdmins)
		protected.GET("/admin/:id", adminCtrl.GetAdminByID)
		protected.PUT("/admin/:id", adminCtrl.UpdateAdmin)
		protected.DELETE("/admin/:id", adminCtrl.DeleteAdmin)

		// Category routes
		protected.POST("/categories", categoryCtrl.CreateCategory)
		protected.GET("/categories", categoryCtrl.GetAllCategories)
		protected.GET("/categories/:id", categoryCtrl.GetCategoryByID)
		protected.PUT("/categories/:id", categoryCtrl.UpdateCategory)
		protected.DELETE("/categories/:id", categoryCtrl.DeleteCategory)

		// Product routes
		protected.POST("/products", productCtrl.CreateProduct)
		protected.GET("/all-products", productCtrl.GetAllProducts)
		protected.GET("/products", productCtrl.GetAllPaginatedProducts)
		protected.GET("/products/:id", productCtrl.GetProductByID)
		protected.PUT("/products/:id", productCtrl.UpdateProduct)
		protected.DELETE("/products/:id", productCtrl.DeleteProduct)

		// Transaction routes
		protected.POST("/transactions", transactionCtrl.CreateTransaction)
		protected.GET("/transactions/history", transactionCtrl.GetAllTransactions)
		protected.GET("/transactions/:id", transactionCtrl.GetTransactionByID)
		protected.GET("/transactions/by-product", transactionCtrl.GetTransactionsByProduct)

		// protected.GET("/transactions/history", transactionCtrl.GetHistory)
	}
}
