package utils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func PaginatedSuccessResponse(c *gin.Context, code int, message string, data interface{}, meta *Meta) {
	c.JSON(code, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorResponse(c *gin.Context, code int, error string) {
	c.JSON(code, Response{
		Success: false,
		Error:   error,
	})
}

func ValidationErrorResponse(c *gin.Context, err validator.ValidationErrors, obj interface{}) {
	errorDetails := make(map[string]string)

	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	for _, e := range err {
		var jsonField string

		// Ambil field dari objek berdasarkan nama field di struct
		if fieldStruct, ok := objType.FieldByName(e.StructField()); ok {
			jsonTag := fieldStruct.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				jsonField = jsonTag
			}
		}

		// Fallback ke nama field Go jika tag `json` tidak ada
		if jsonField == "" {
			jsonField = e.Field()
		}

		// Buat pesan berdasarkan tag validasi
		var message string
		switch e.Tag() {
		case "transaction_type":
			message = "wrong transaction_type, choose either 'stock_out' or 'stock_in'"
		case "gender":
			message = "invalid gender, choose either 'L' or 'P'"
		case "date":
			message = "invalid date format, use 'YYYY-MM-DD'"
		case "required":
			message = fmt.Sprintf("%s is required", jsonField)
		default:
			message = fmt.Sprintf("%s is invalid", jsonField)
		}

		errorDetails[jsonField] = message
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Validation failed",
		"details": errorDetails,
	})
}
