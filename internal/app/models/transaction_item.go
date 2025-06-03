package models

import "time"

type TransactionItem struct {
	ID            int       `json:"id,omitempty" db:"id"`
	TransactionID int       `json:"transaksi_id,omitempty" db:"transaksi_id"`
	ProductName   string    `json:"nama_produk,omitempty" db:"nama_produk"`
	ProductID     int       `json:"produk_id" db:"produk_id" validate:"required"`
	Quantity      int       `json:"jumlah" db:"jumlah" validate:"required,min=1"`
	StockBefore   int       `json:"stok_sebelum,omitempty" db:"stok_sebelum"`
	StockAfter    int       `json:"stok_sesudah,omitempty" db:"stok_sesudah"`
	CreatedAt     time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitzero" db:"updated_at"`
}

type CreateTransactionItemRequest struct {
	ProductID int `json:"produk_id" validate:"required"`
	Quantity  int `json:"jumlah" validate:"required,min=1"`
}

type TransactionItemWithProduct struct {
	TransactionItem
	ProductName string `json:"nama_produk" db:"nama_produk"`
}
