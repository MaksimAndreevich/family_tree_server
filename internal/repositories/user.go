package repositories

import (
	"gitlab.com/geno-tree/go-back/internal/database"
	"gitlab.com/geno-tree/go-back/internal/models"
)

type UserRepository struct {
	db *database.Db
}

func NewUserRepository(db *database.Db) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	result := repo.db.Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := repo.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
