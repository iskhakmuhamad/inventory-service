package models

import (
	"time"
)

type Transaction struct {
	ID        int               `json:"id" db:"id"`
	Type      string            `json:"jenis_transaksi" db:"jenis_transaksi" validate:"required,transaction_type"`
	Date      time.Time         `json:"tanggal_transaksi,omitzero" db:"tanggal_transaksi"`
	CreatedBy int               `json:"created_by,omitempty" db:"created_by"`
	Notes     string            `json:"keterangan,omitempty" db:"keterangan"`
	CreatedAt time.Time         `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at,omitzero" db:"updated_at"`
	Items     []TransactionItem `json:"items" validate:"required,min=1,dive"`
}

type CreateTransactionRequest struct {
	Type  string                         `json:"jenis_transaksi" validate:"required,transaction_type"`
	Notes string                         `json:"keterangan"`
	Items []CreateTransactionItemRequest `json:"items" validate:"required,min=1,dive"`
}

type TransactionWithDetails struct {
	Transaction
	AdminName string `json:"admin_name" db:"admin_name"`
}
