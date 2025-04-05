package services

import (
	"dialogue/internal/models"
	"dialogue/internal/store"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Authenticate(string, string) (*models.User, error)
	PasswordUpdate(int, string, string) (err error)
}

type UserStruct struct {
	db store.Store
}

// Authenticate authenticates a user with given data.
func (us *UserStruct) Authenticate(email, password string) (*models.User, error) {
	user, _ := us.db.User().FindByEmail(email)

	// Check whether the hashed password and plain-text password provided match. If correct, return user ID.
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (us *UserStruct) PasswordUpdate(id int, currentPassword, newPassword string) (err error) {
	var user *models.User
	user, err = us.db.User().FindByID(id)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(currentPassword))
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

	err = us.db.User().PasswordUpdate(id, newHashedPassword)
	if err != nil {
		return err
	}
	return nil
}
