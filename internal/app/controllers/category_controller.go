package controllers

import (
	"net/http"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/services"
	"github.com/iskhakmuhamad/inventory-service/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryController struct {
	Service  *services.CategoryService
	Validate *validator.Validate
}

func NewCategoryController(service *services.CategoryService, validate *validator.Validate) *CategoryController {
	return &CategoryController{Service: service, Validate: validate}
}

func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	var category models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newCategory := models.Category{
		Name:        category.Name,
		Description: category.Description,
	}

	err := ctrl.Validate.Struct(newCategory)
	if err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), category)
		return
	}

	if err := ctrl.Service.CreateCategory(&newCategory); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category created successfully", newCategory)
}

func (ctrl *CategoryController) GetCategoryByID(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	category, err := ctrl.Service.GetCategoryByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

func (ctrl *CategoryController) GetAllCategories(c *gin.Context) {
	categories, err := ctrl.Service.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	var category models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), category)
		return
	}

	// Map request to model
	updatedCategory := models.Category{
		ID:          id,
		Name:        category.Name,
		Description: category.Description,
	}

	if err := ctrl.Service.UpdateCategory(&updatedCategory); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", updatedCategory)
}

func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	if err := ctrl.Service.DeleteCategory(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
