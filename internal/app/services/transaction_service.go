package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/repositories"
	"github.com/iskhakmuhamad/inventory-service/pkg/constants"
	"github.com/jmoiron/sqlx"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
	db   *sqlx.DB
}

func NewTransactionService(repo *repositories.TransactionRepository, db *sqlx.DB) *TransactionService {
	return &TransactionService{repo: repo, db: db}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	user, ok := ctx.Value(constants.CtxUserKey).(models.AuthenticatedUser)
	if !ok || user.UserID == "" {
		return fmt.Errorf("[InventoryService] [TransactionService] [CreateTransaction] [GotUserToken]: Doent Get User Data ")
	}

	transaction.CreatedBy, _ = strconv.Atoi(user.UserID)

	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	for _, item := range transaction.Items {
		if transaction.Type == "stock_out" && item.Quantity <= 0 {
			return fmt.Errorf("quantity for stock out must be greater than zero")
		}
	}

	// Create the transaction (either Stock In or Stock Out)
	if err := s.repo.Create(tx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit the transaction (all operations are successful)
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Get all transactions
func (s *TransactionService) GetAllTransactions() ([]models.TransactionWithDetails, error) {
	return s.repo.GetAll()
}

// Get transaction details by ID
func (s *TransactionService) GetTransactionDetails(transactionID int) (*models.TransactionWithDetails, error) {
	return s.repo.GetByID(transactionID)
}

func (s *TransactionService) GetTransactionsByProduct(productID string) ([]models.TransactionWithDetails, error) {
	return s.repo.GetTransactionsByProduct(productID)
}
