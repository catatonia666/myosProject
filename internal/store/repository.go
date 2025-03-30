package store

import "dialogue/internal/models"

type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	Get(int) (*models.User, error)
	PasswordUpdate(int, string, string) error
	Exists(int) (bool, error)
}

type StoryRepository interface {
	Create(any) (int, error)
	Get(string, int) (any, error)
	Update(int, any) error
	Delete(any) error
	DeleteWholeStory(int) error

	StoriesToDisplay(int, int) ([]models.StartingBlock, error)
	CreatedBlocks(int) ([]models.CommonBlock, error)
	WholeStory(int) (models.StartingBlock, []models.CommonBlock, error)

	RetrieveBlocks(id int) (retrievedBlocks models.RelatedToStoryBlocks, err error)
}
