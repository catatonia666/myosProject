package sqlstore

import (
	"dialogue/internal/models"
	"errors"
	"fmt"

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

func (ur *UserRepository) Find(int) (*models.User, error) {
	return nil, nil
}

func (ur *UserRepository) FindByEmail(string) (*models.User, error) {
	return nil, nil
}

// Authenticate authenticates a user with given data.
func (ur *UserRepository) Authenticate(email, password string) (*models.User, error) {
	var user models.User

	// Retrieve ID and hashed password associated with the given email. If no matching email exists then return an error.
	result := ur.store.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, result.Error
		}
	}

	fmt.Printf("User from Authenticate: %+v\n", user)
	fmt.Printf("User ID Type: %T\n", user.ID)

	// Check whether the hashed password and plain-text password provided match. If correct, return user ID.
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}
	return &user, nil
}

// Exists checks if user with provided ID exists in the database.
func (ur *UserRepository) Exists(id int) (bool, error) {
	var exists bool
	err := ur.store.db.Table("users").Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error
	return exists, err
}

// GetUser gets user with provided ID if exists.
func (ur *UserRepository) GetUser(id int) (*models.User, error) {
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
