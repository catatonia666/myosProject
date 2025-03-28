package services

import (
	"dialogue/internal/models"
	"dialogue/internal/store"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Authenticate(string, string) (*models.User, error)
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
