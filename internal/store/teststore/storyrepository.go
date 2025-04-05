package teststore

import "dialogue/internal/models"

type StoryRepository struct {
}

func (sr *StoryRepository) Create(any) (int, error) {
	return 0, nil
}

func (sr *StoryRepository) Get(string, int) (any, error) {
	return nil, nil
}

func (sr *StoryRepository) Update(int, any) error {
	return nil
}

func (sr *StoryRepository) Delete(any) error {
	return nil
}

func (sr *StoryRepository) DeleteWholeStory(int) error {
	return nil
}

func (sr *StoryRepository) StoriesToDisplay(int, int) ([]models.StartingBlock, error) {
	return nil, nil
}

func (sr *StoryRepository) WholeStory(int) (models.StartingBlock, []models.CommonBlock, models.RelatedToStoryBlocks, error) {
	return models.StartingBlock{}, nil, models.RelatedToStoryBlocks{}, nil
}

func (sr *StoryRepository) GetAllStories() (stories []models.StartingBlock, err error) {
	return nil, nil
}

func (sr *StoryRepository) CreatedBlocks(int) ([]models.CommonBlock, error) {
	return nil, nil
}
