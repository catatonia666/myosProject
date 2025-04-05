package sqlstore

import (
	"dialogue/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	store *Store
}

// Create inserts a new user into the database.
func (ur *UserRepository) Create(u *models.User) error {
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if result := ur.store.db.Create(u); result.Error != nil {
		return models.ErrDuplicateEmail
	}
	return nil
}

// FindByID gets user with provided ID if exists.
func (ur *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := ur.store.db.Table("users").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &user, nil
}

// FindByEmail gets user with provided email if exists.
func (ur *UserRepository) FindByEmail(email string) (user *models.User, err error) {
	result := ur.store.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, result.Error
		}
	}

	return user, nil
}

// PasswordUpdate updates user's password.
func (ur *UserRepository) PasswordUpdate(id int, newHashedPassword []byte) (err error) {
	err = ur.store.db.Table("users").Where("id = ?", id).Update("hashed_password", &newHashedPassword).Error
	if err != nil {
		return err
	}
	return nil
}
