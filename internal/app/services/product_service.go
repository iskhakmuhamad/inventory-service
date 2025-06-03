package services

import (
	"math"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/repositories"
	"github.com/iskhakmuhamad/inventory-service/pkg/utils"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

type ProductServiceInterface interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id int) (*models.ProductWithCategory, error)
	GetAllProducts() ([]models.ProductWithCategory, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id int) error
	GetPaginatedProductsWithMeta(page, limit int) ([]models.ProductWithCategory, *utils.Meta, error)
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.Repo.Create(product)
}

func (s *ProductService) GetProductByID(id int) (*models.ProductWithCategory, error) {
	return s.Repo.GetByID(id)
}

func (s *ProductService) GetAllProducts() ([]models.ProductWithCategory, error) {
	return s.Repo.GetAll()
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.Repo.Update(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.Repo.Delete(id)
}

func (s *ProductService) GetPaginatedProductsWithMeta(page, limit int) ([]models.ProductWithCategory, *utils.Meta, error) {
	offset := (page - 1) * limit

	products, err := s.Repo.GetPaginated(offset, limit)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.Repo.CountAll()
	if err != nil {
		return nil, nil, err
	}

	meta := &utils.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return products, meta, nil
}
