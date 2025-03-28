package sqlstore

import (
	"dialogue/internal/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	store *Store
}

// Insert insets a new user into the database.
func (ur *UserRepository) Create(u *models.User) error {
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if result := ur.store.db.Create(u); result.Error != nil {
		return models.ErrDuplicateEmail
	}
	return nil
}

// GetUser gets user with provided ID if exists.
func (ur *UserRepository) Get(id int) (*models.User, error) {
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

func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user *models.User
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

// Exists checks if user with provided ID exists in the database.
func (ur *UserRepository) Exists(id int) (bool, error) {
	var exists bool
	err := ur.store.db.Table("users").Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error
	return exists, err
}

// PasswordUpdate updates user's password.
func (ur *UserRepository) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var pw models.User
	err := ur.store.db.Table("users").Select("hashed_password").Where("id = ?", id).Scan(&pw).Error
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pw.HashedPassword), []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrInvalidCredentials
		} else {
			return err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	err = ur.store.db.Table("users").Where("id = ?", id).Update("hashed_password", &newHashedPassword).Error
	return err
}
