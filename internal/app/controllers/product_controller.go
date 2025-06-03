package controllers

import (
	"net/http"
	"strconv"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/services"
	"github.com/iskhakmuhamad/inventory-service/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	Service  services.ProductServiceInterface
	Validate *validator.Validate
}

func NewProductController(service services.ProductServiceInterface, validate *validator.Validate) *ProductController {
	return &ProductController{Service: service, Validate: validate}
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := ctrl.Validate.Struct(product)
	if err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), product)
		return
	}

	if err := ctrl.Service.CreateProduct(&product); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product created successfully", product)
}

func (ctrl *ProductController) GetProductByID(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	product, err := ctrl.Service.GetProductByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.Service.GetAllProducts()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}

func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), product)
		return
	}

	product.ID = id
	if err := ctrl.Service.UpdateProduct(&product); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	if err := ctrl.Service.DeleteProduct(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

func (ctrl *ProductController) GetAllPaginatedProducts(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr != "" && limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		products, meta, err := ctrl.Service.GetPaginatedProductsWithMeta(page, limit)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PaginatedSuccessResponse(c, http.StatusOK, "Products retrieved successfully", products, meta)
		return
	}

	products, err := ctrl.Service.GetAllProducts()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}
