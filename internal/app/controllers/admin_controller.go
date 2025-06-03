package controllers

import (
	"log"
	"net/http"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/services"
	"github.com/iskhakmuhamad/inventory-service/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AdminController struct {
	Service  *services.AdminService
	Validate *validator.Validate
}

func NewAdminController(service *services.AdminService, validate *validator.Validate) *AdminController {
	return &AdminController{
		Service:  service,
		Validate: validate,
	}
}

func (ctrl *AdminController) CreateAdmin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		log.Println("[InventoryService] [Err] [AdminController] [CreateAdmin] [ShouldbindJSON] : ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := ctrl.Validate.Struct(admin)
	if err != nil {
		log.Println("[InventoryService] [Err] [AdminController] [CreateAdmin] [Validation] : ", err)
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), admin)
		return
	}

	if err := ctrl.Service.CreateAdmin(&admin); err != nil {
		log.Println("[InventoryService] [Err] [AdminController] [CreateAdmin] [Service] : ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin created successfully", admin)
}

func (ctrl *AdminController) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println("[InventoryService] [Err] [AdminController] [Login] [ShouldBindJSON] : ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctrl.Validate.Struct(loginData); err != nil {
		log.Println("[InventoryService] [Err] [AdminController] [Login] [Validate] : ", err)
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), loginData)
		return
	}

	token, err := ctrl.Service.Login(loginData.Email, loginData.Password)
	if err != nil {
		log.Println("[InventoryService] [Err] [AdminController] [Login] [Service] : ", err)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

func (ctrl *AdminController) GetAdminByID(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	admin, err := ctrl.Service.GetAdminByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Admin not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin retrieved successfully", admin)
}

func (ctrl *AdminController) GetProfile(c *gin.Context) {
	admin, err := ctrl.Service.GetProfile(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Profile not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", admin)
}

func (ctrl *AdminController) GetAllAdmins(c *gin.Context) {
	admins, err := ctrl.Service.GetAllAdmins()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Admins retrieved successfully", admins)
}

func (ctrl *AdminController) UpdateAdmin(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}

	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), admin)
		return
	}

	admin.ID = id
	if err := ctrl.Service.UpdateAdmin(&admin); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin updated successfully", admin)
}

func (ctrl *AdminController) UpdateProfile(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		utils.ValidationErrorResponse(c, err.(validator.ValidationErrors), admin)
		return
	}

	if err := ctrl.Service.UpdateProfile(c.Request.Context(), &admin); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", admin)
}

func (ctrl *AdminController) DeleteAdmin(c *gin.Context) {
	id, ok := utils.GetIntParam(c, "id")
	if !ok {
		return
	}
	if err := ctrl.Service.DeleteAdmin(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin deleted successfully", nil)
}
