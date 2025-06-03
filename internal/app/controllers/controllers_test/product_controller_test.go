package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/iskhakmuhamad/inventory-service/internal/app/controllers"
	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	mockService "github.com/iskhakmuhamad/inventory-service/internal/app/services/mock"
)

// --- UNIT TEST ---
func TestGetAllProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mockService.ProductServiceInterface)
	controller := &controllers.ProductController{
		Service: mockService,
	}

	t.Run("success without pagination", func(t *testing.T) {
		expectedProducts := []models.ProductWithCategory{{
			Product: models.Product{
				ID:   1,
				Name: "Iphone",
			},
			CategoryName: "Electronic",
		},
		}

		mockService.On("GetAllProducts").Return(expectedProducts, nil)

		req, _ := http.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		controller.GetAllProducts(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Products retrieved successfully")
		mockService.AssertExpectations(t)
	})
}
