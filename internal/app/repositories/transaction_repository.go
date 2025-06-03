package repositories

import (
	"fmt"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create a transaction (either Stock In or Stock Out)
func (r *TransactionRepository) Create(tx *sqlx.Tx, transaction *models.Transaction) error {
	// Insert into transactions table
	query := `INSERT INTO transactions (jenis_transaksi, created_by, keterangan) VALUES (?, ?, ?)`

	result, err := tx.Exec(query, transaction.Type, transaction.CreatedBy, transaction.Notes)
	if err != nil {
		return err
	}

	// Get the generated transaction ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	transaction.ID = int(id)

	// Create transaction items (related products in the transaction)
	for _, item := range transaction.Items {
		if err := r.CreateItem(tx, &item, transaction.ID, transaction.Type); err != nil {
			return err
		}
	}

	return nil
}

// Create transaction items (product-related in a transaction)
func (r *TransactionRepository) CreateItem(tx *sqlx.Tx, item *models.TransactionItem, transactionID int, transactionType string) error {
	// Get the product's current stock
	var currentStock int
	err := tx.Get(&currentStock, "SELECT stok_produk FROM products WHERE id = ?", item.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Validate Stock Out (ensure the transaction does not exceed available stock)
	if transactionType == "stock_out" && item.Quantity > currentStock {
		return fmt.Errorf("not enough stock for product ID %d", item.ProductID)
	}

	// Adjust stock based on transaction type
	var newStock int
	if transactionType == "stock_in" {
		newStock = currentStock + item.Quantity
	} else if transactionType == "stock_out" {
		newStock = currentStock - item.Quantity
	}

	// Update stock in the products table
	_, err = tx.Exec("UPDATE products SET stok_produk = ? WHERE id = ?", newStock, item.ProductID)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Create the transaction item record in the database
	query := `INSERT INTO transaction_items (transaksi_id, produk_id, jumlah, stok_sebelum, stok_sesudah) 
			  VALUES (?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, transactionID, item.ProductID, item.Quantity, currentStock, newStock)
	return err
}

func (r *TransactionRepository) GetByID(id int) (*models.TransactionWithDetails, error) {
	var transaction models.TransactionWithDetails
	query := `SELECT t.*, CONCAT(a.nama_depan, ' ', a.nama_belakang) AS admin_name 
			  FROM transactions t 
			  LEFT JOIN admins a ON t.created_by = a.id 
			  WHERE t.id = ?`
	err := r.db.Get(&transaction, query, id)
	if err != nil {
		return nil, err
	}

	// Get transaction items
	items, err := r.GetTransactionItems(id)
	if err != nil {
		return nil, err
	}

	transaction.Items = items
	return &transaction, nil
}

// Get all transaction items
func (r *TransactionRepository) GetTransactionItems(transactionID int) ([]models.TransactionItem, error) {
	var items []models.TransactionItem
	query := `
		SELECT 
			ti.*, 
			p.nama_produk AS nama_produk
		FROM transaction_items ti
		JOIN products p ON ti.produk_id = p.id
		WHERE ti.transaksi_id = ?
	`
	err := r.db.Select(&items, query, transactionID)
	return items, err
}

func (r *TransactionRepository) GetTransactionByProductItems(transactionID int, productID string) ([]models.TransactionItem, error) {
	var items []models.TransactionItem
	query := `SELECT * FROM transaction_items WHERE transaksi_id = ? AND produk_id = ? `
	err := r.db.Select(&items, query, transactionID, productID)
	return items, err
}

func (r *TransactionRepository) GetAll() ([]models.TransactionWithDetails, error) {
	var transactions []models.TransactionWithDetails

	// Query to get all transactions along with admin name (for better context)
	query := `SELECT t.*, CONCAT(a.nama_depan, ' ', a.nama_belakang) AS admin_name 
			  FROM transactions t 
			  LEFT JOIN admins a ON t.created_by = a.id 
			  ORDER BY t.created_at DESC`

	// Fetch the transactions from the database
	err := r.db.Select(&transactions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Get the transaction items for each transaction
	for i, transaction := range transactions {
		// Query to get items related to the transaction
		items, err := r.GetTransactionItems(transaction.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction items for transaction %d: %w", transaction.ID, err)
		}
		transactions[i].Items = items
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByProduct(productID string) ([]models.TransactionWithDetails, error) {
	var transactions []models.TransactionWithDetails

	query := `
		SELECT DISTINCT t.*, CONCAT(a.nama_depan, ' ', a.nama_belakang) AS admin_name 
		FROM transactions t
		LEFT JOIN admins a ON t.created_by = a.id 
		JOIN transaction_items ti ON t.id = ti.transaksi_id
		WHERE ti.produk_id = ?
		ORDER BY t.created_at DESC
	`

	err := r.db.Select(&transactions, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	for i, transaction := range transactions {
		items, err := r.GetTransactionByProductItems(transaction.ID, productID)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction items for transaction %d: %w", transaction.ID, err)
		}
		transactions[i].Items = items
	}

	return transactions, nil
}
