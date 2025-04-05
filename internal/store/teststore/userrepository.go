package teststore

import (
	"dialogue/internal/models"
)

type UserRepository struct {
	store *Store
	users map[int]*models.User
}

func (ur *UserRepository) Create(u *models.User) error {
	for _, existingUser := range ur.users {
		if existingUser.Email == u.Email {
			return models.ErrDuplicateEmail
		}
	}

	u.ID = len(ur.users) + 1
	ur.users[u.ID] = u

	return nil
}

func (ur *UserRepository) FindByID(id int) (*models.User, error) {

	u, ok := ur.users[id]
	if !ok {
		return nil, models.ErrNoRecord
	}

	return u, nil
}

func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	for _, u := range ur.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, models.ErrNoRecord
}

func (ur *UserRepository) Authenticate(string, string) (*models.User, error) {
	return nil, nil
}

func (ur *UserRepository) PasswordUpdate(int, []byte) error {
	return nil
}
