package models

import (
	"log"
	"time"
)

type Admin struct {
	ID           int       `json:"id" db:"id"`
	NamaDepan    string    `json:"nama_depan" db:"nama_depan" validate:"required,min=2,max=100"`
	NamaBelakang string    `json:"nama_belakang" db:"nama_belakang" validate:"required,min=2,max=100"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	TanggalLahir string    `json:"tanggal_lahir" db:"tanggal_lahir" validate:"required,date"`
	JenisKelamin string    `json:"jenis_kelamin" db:"jenis_kelamin" validate:"required,gender"`
	Password     string    `json:"password,omitempty" db:"password" validate:"required,min=6"`
	CreatedAt    time.Time `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitzero" db:"updated_at"`
}

type AdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminResponse struct {
	ID           int       `json:"id"`
	NamaDepan    string    `json:"nama_depan"`
	NamaBelakang string    `json:"nama_belakang"`
	Email        string    `json:"email"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	CreatedAt    time.Time `json:"created_at,omitzero"`
	UpdatedAt    time.Time `json:"updated_at,omitzero"`
}

func (a *Admin) ToResponse() *AdminResponse {
	tanggalLahirParsed, err := time.Parse("2006-01-02", a.TanggalLahir)
	if err != nil {
		log.Println("[InventoryService] [Err] [Admin] [Model] [ToResponse]: ", err)
		return nil
	}
	return &AdminResponse{
		ID:           a.ID,
		NamaDepan:    a.NamaDepan,
		NamaBelakang: a.NamaBelakang,
		Email:        a.Email,
		TanggalLahir: tanggalLahirParsed,
		JenisKelamin: a.JenisKelamin,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

type AuthenticatedUser struct {
	UserID string
	Email  string
}
