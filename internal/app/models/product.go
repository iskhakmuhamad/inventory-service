package models

import (
	"time"
)

type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"nama_produk" db:"nama_produk" validate:"required,min=2,max=255"`
	Description string    `json:"deskripsi_produk" db:"deskripsi_produk"`
	Image       string    `json:"gambar_produk" db:"gambar_produk"`
	CategoryID  int       `json:"kategori_produk_id" db:"kategori_produk_id" validate:"required"`
	Stock       int       `json:"stok_produk" db:"stok_produk" validate:"min=0"`
	CreatedAt   time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitzero" db:"updated_at"`
}

type ProductWithCategory struct {
	Product
	CategoryName string `json:"nama_kategori" db:"nama_kategori"`
}

type CreateProductRequest struct {
	Name        string `json:"nama_produk" validate:"required,min=2,max=255"`
	Description string `json:"deskripsi_produk"`
	Image       string `json:"gambar_produk"`
	CategoryID  int    `json:"kategori_produk_id" validate:"required"`
	Stock       int    `json:"stok_produk" validate:"min=0"`
}

type UpdateProductRequest struct {
	Name        string `json:"nama_produk" validate:"required,min=2,max=255"`
	Description string `json:"deskripsi_produk"`
	Image       string `json:"gambar_produk"`
	CategoryID  int    `json:"kategori_produk_id" validate:"required"`
	Stock       int    `json:"stok_produk" validate:"min=0"`
}
