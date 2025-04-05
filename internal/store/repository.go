package store

import "dialogue/internal/models"

type UserRepository interface {
	Create(*models.User) error
	FindByID(int) (*models.User, error)
	FindByEmail(string) (*models.User, error)
	PasswordUpdate(int, []byte) error
}

type StoryRepository interface {
	Create(any) (int, error)
	Get(string, int) (any, error)
	Update(int, any) error
	Delete(any) error
	DeleteWholeStory(int) error

	StoriesToDisplay(int, int) ([]models.StartingBlock, error)
	WholeStory(int) (models.StartingBlock, []models.CommonBlock, models.RelatedToStoryBlocks, error)
	GetAllStories() ([]models.StartingBlock, error)

	CreatedBlocks(int) ([]models.CommonBlock, error)
}
