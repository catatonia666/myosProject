package store

import "dialogue/internal/models"

type UserRepository interface {
	Create(*models.User) error
	Find(int) (*models.User, error)
	FindByEmail(string) (*models.User, error)
	Authenticate(string, string) (*models.User, error)
	GetUser(int) (*models.User, error)
	PasswordUpdate(int, string, string) error
	Exists(int) (bool, error)
}

type StoryRepository interface {
	Latest(int) []models.FirstBlock
	CreateFB(int, string, string, []string, bool) int
	CreatedFBView(int) models.DialoguesData
	EditFB(int, int, string, string, []string)
	DeleteFB(int)
	EditBView(int) models.DialoguesData
	EditB(int, int, string, string, []string)
	DeleteB(int)
}
