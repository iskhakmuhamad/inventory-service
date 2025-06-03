package repositories

import (
	"github.com/iskhakmuhamad/inventory-service/internal/app/models"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (nama_kategori, deskripsi_kategori) VALUES (?, ?)`

	result, err := r.db.Exec(query, category.Name, category.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	category.ID = int(id)
	return nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE id = ?`
	err := r.db.Get(&category, query, id)
	return &category, err
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	query := `SELECT * FROM categories ORDER BY created_at DESC`
	err := r.db.Select(&categories, query)
	return categories, err
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := `UPDATE categories SET nama_kategori = ?, deskripsi_kategori = ? WHERE id = ?`
	_, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	return err
}

func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
