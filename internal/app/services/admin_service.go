package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iskhakmuhamad/inventory-service/internal/app/models"
	"github.com/iskhakmuhamad/inventory-service/internal/app/repositories"
	"github.com/iskhakmuhamad/inventory-service/pkg/constants"
)

type AdminService struct {
	Repo      *repositories.AdminRepository
	JwtSecret string
}

func NewAdminService(repo *repositories.AdminRepository, jwtSecret string) *AdminService {
	return &AdminService{Repo: repo, JwtSecret: jwtSecret}
}

func (s *AdminService) CreateAdmin(admin *models.Admin) error {
	return s.Repo.Create(admin)
}

func (s *AdminService) GetAdminByID(id int) (*models.Admin, error) {
	return s.Repo.GetByID(id)
}

func (s *AdminService) GetAllAdmins() ([]models.Admin, error) {
	return s.Repo.GetAll()
}

func (s *AdminService) UpdateAdmin(admin *models.Admin) error {
	return s.Repo.Update(admin)
}

func (s *AdminService) UpdateProfile(ctx context.Context, admin *models.Admin) error {
	user, ok := ctx.Value(constants.CtxUserKey).(models.AuthenticatedUser)
	if !ok || user.UserID == "" {
		return fmt.Errorf("[InventoryService] [AdminService] [UpdateProfile] [GotUserToken]: Doesnt Get User Data ")
	}

	admin.ID, _ = strconv.Atoi(user.UserID)

	return s.Repo.Update(admin)
}

func (s *AdminService) GetProfile(ctx context.Context) (*models.Admin, error) {
	user, ok := ctx.Value(constants.CtxUserKey).(models.AuthenticatedUser)
	if !ok || user.UserID == "" {
		return nil, fmt.Errorf("[InventoryService] [AdminService] [GetProfile] [GotUserToken]: Doesnt Get User Data ")
	}
	ID, _ := strconv.Atoi(user.UserID)

	return s.Repo.GetByID(ID)
}

func (s *AdminService) DeleteAdmin(id int) error {
	return s.Repo.Delete(id)
}

func (s *AdminService) Login(email, password string) (string, error) {
	admin, err := s.Repo.FindAdminByEmail(email)
	if err != nil {
		log.Println("[InventoryService] [Err] [AdminService] [Login] [FindAdminByEmail]: ", err)
		return "", err
	}

	if admin.Password != password {
		log.Println("[InventoryService] [Err] [AdminService] [Password]: Incorect User Password")
		return "", errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"email": admin.Email,
		"sub":   strconv.Itoa(admin.ID),
		"iat":   time.Now().Unix(),                     // Issued at time
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expiration time (1 day)
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
