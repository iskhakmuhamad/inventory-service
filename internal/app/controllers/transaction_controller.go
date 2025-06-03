package controllers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/services"
	"github.com/iskhakmuhamad/inventory-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	Service  *services.TransactionService
	Validate *validator.Validate
}

func NewTransactionController(service *services.TransactionService, validate *validator.Validate) *TransactionController {
	return &TransactionController{Service: service, Validate: validate}
}

func (ctrl *TransactionController) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err := ctrl.Validate.Struct(transaction)
	if err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), transaction)
		return
	}

	// Call service to create transaction
	if err := ctrl.Service.CreateTransaction(c.Request.Context(), &transaction); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction created successfully", transaction)
}

// Get all transactions
func (ctrl *TransactionController) GetAllTransactions(c *gin.Context) {
	transactions, err := ctrl.Service.GetAllTransactions()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// Get a transaction by ID
func (ctrl *TransactionController) GetTransactionByID(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	transaction, err := ctrl.Service.GetTransactionDetails(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

func (ctrl *TransactionController) GetTransactionsByProduct(c *gin.Context) {
	productID := c.Query("product_id")
	if productID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "product_id is required")
		return
	}

	transactions, err := ctrl.Service.GetTransactionsByProduct(productID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}
