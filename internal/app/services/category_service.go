package services

import (
	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/repositories"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{Repo: repo}
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.Repo.Create(category)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.Repo.GetByID(id)
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.Repo.GetAll()
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.Repo.Update(category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.Repo.Delete(id)
}
