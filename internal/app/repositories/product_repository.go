package repositories

import (
	"github.com/iskhakmuhamad/inventory-service/internal/app/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := `INSERT INTO products (nama_produk, deskripsi_produk, gambar_produk, kategori_produk_id, stok_produk) 
			  VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, product.Name, product.Description, product.Image,
		product.CategoryID, product.Stock)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = int(id)
	return nil
}

func (r *ProductRepository) GetByID(id int) (*models.ProductWithCategory, error) {
	var product models.ProductWithCategory
	query := `SELECT p.*, c.nama_kategori FROM products p 
			  LEFT JOIN categories c ON p.kategori_produk_id = c.id 
			  WHERE p.id = ?`
	err := r.db.Get(&product, query, id)
	return &product, err
}

func (r *ProductRepository) GetAll() ([]models.ProductWithCategory, error) {
	var products []models.ProductWithCategory
	query := `SELECT p.*, c.nama_kategori FROM products p 
			  LEFT JOIN categories c ON p.kategori_produk_id = c.id 
			  ORDER BY p.created_at DESC`
	err := r.db.Select(&products, query)
	return products, err
}

func (r *ProductRepository) GetPaginated(offset, limit int) ([]models.ProductWithCategory, error) {
	var products []models.ProductWithCategory

	query := `
		SELECT p.*, c.nama_kategori 
		FROM products p
		LEFT JOIN categories c ON p.kategori_produk_id = c.id
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.Select(&products, query, limit, offset)
	return products, err
}

func (r *ProductRepository) CountAll() (int, error) {
	var total int
	err := r.db.Get(&total, "SELECT COUNT(*) FROM products")
	return total, err
}

func (r *ProductRepository) Update(product *models.Product) error {
	query := `UPDATE products SET nama_produk = ?, deskripsi_produk = ?, gambar_produk = ?, 
			  kategori_produk_id = ?, stok_produk = ? WHERE id = ?`

	_, err := r.db.Exec(query, product.Name, product.Description, product.Image,
		product.CategoryID, product.Stock, product.ID)
	return err
}

func (r *ProductRepository) UpdateStock(id int, stock int) error {
	query := `UPDATE products SET stok_produk = ? WHERE id = ?`
	_, err := r.db.Exec(query, stock, id)
	return err
}

func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ProductRepository) GetStock(id int) (int, error) {
	var stock int
	query := `SELECT stok_produk FROM products WHERE id = ?`
	err := r.db.Get(&stock, query, id)
	return stock, err
}
