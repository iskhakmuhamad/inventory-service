package models

import (
	"time"
)

type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"nama_kategori" db:"nama_kategori" validate:"required,min=2,max=255"`
	Description string    `json:"deskripsi_kategori" db:"deskripsi_kategori"`
	CreatedAt   time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitzero" db:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"nama_kategori" validate:"required,min=2,max=255"`
	Description string `json:"deskripsi_kategori"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"nama_kategori" validate:"required,min=2,max=255"`
	Description string `json:"deskripsi_kategori"`
}
