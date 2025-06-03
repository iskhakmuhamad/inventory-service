package repositories

import (
	"log"

	"github.com/iskhakmuhamad/inventory-service/internal/app/models"

	"github.com/jmoiron/sqlx"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(admin *models.Admin) error {
	query := `INSERT INTO admins (nama_depan, nama_belakang, email, tanggal_lahir, jenis_kelamin, password) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, admin.NamaDepan, admin.NamaBelakang, admin.Email,
		admin.TanggalLahir, admin.JenisKelamin, admin.Password)
	if err != nil {
		log.Println("[InventoryService] [Err] [AdminRepository] [Create] [Exec] : ", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("[InventoryService] [Err] [AdminRepository] [Create] [GotLastInsertedID] : ", err)
		return err
	}

	admin.ID = int(id)
	return nil
}

func (r *AdminRepository) GetByID(id int) (*models.Admin, error) {
	var admin models.Admin
	query := `SELECT * FROM admins WHERE id = ?`
	err := r.db.Get(&admin, query, id)
	return &admin, err
}

func (r *AdminRepository) GetByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	query := `SELECT * FROM admins WHERE email = ?`
	err := r.db.Get(&admin, query, email)
	return &admin, err
}

func (r *AdminRepository) GetAll() ([]models.Admin, error) {
	var admins []models.Admin
	query := `SELECT * FROM admins ORDER BY created_at DESC`
	err := r.db.Select(&admins, query)
	return admins, err
}

func (r *AdminRepository) Update(admin *models.Admin) error {
	query := `UPDATE admins SET nama_depan = ?, nama_belakang = ?, email = ?, 
			  tanggal_lahir = ?, jenis_kelamin = ?, password = ? WHERE id = ?`

	_, err := r.db.Exec(query, admin.NamaDepan, admin.NamaBelakang, admin.Email,
		admin.TanggalLahir, admin.JenisKelamin, admin.Password, admin.ID)
	return err
}

func (r *AdminRepository) Delete(id int) error {
	query := `DELETE FROM admins WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AdminRepository) FindAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	query := "SELECT id, nama_depan, nama_belakang, email, password FROM admins WHERE email = ?"
	err := r.db.Get(&admin, query, email)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
