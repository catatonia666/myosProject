package teststore

import "dialogue/internal/models"

type StoryRepository struct {
}

func (sr *StoryRepository) RetrieveBlocks(id int) (retrievedBlocks models.RelatedToStoryBlocks, err error) {
	return models.RelatedToStoryBlocks{}, err
}

func (sr *StoryRepository) StoriesToDisplay(int, int) ([]models.StartingBlock, error) {
	return nil, nil
}

func (sr *StoryRepository) Delete(any) error {
	return nil
}

func (sr *StoryRepository) WholeStory(int) (models.StartingBlock, []models.CommonBlock, error) {
	return models.StartingBlock{}, nil, nil
}

func (sr *StoryRepository) CreatedBlocks(int) ([]models.CommonBlock, error) {
	return nil, nil
}

func (sr *StoryRepository) Create(any) (int, error) {
	return 0, nil
}

func (sr *StoryRepository) CreatedFBView(int) models.DialoguesData {
	return models.DialoguesData{}
}

func (sr *StoryRepository) DeleteWholeStory(int) error {
	return nil
}

func (sr *StoryRepository) EditBView(int) models.DialoguesData {
	return models.DialoguesData{}
}

func (sr *StoryRepository) Update(id int, data any) error {
	return nil
}

func (sr *StoryRepository) Get(model string, id int) (any, error) {
	return nil, nil
}
