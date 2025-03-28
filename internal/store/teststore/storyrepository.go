package teststore

import "dialogue/internal/models"

type StoryRepository struct {
}

func (sr *StoryRepository) StoriesToDisplay(int, int) ([]models.StartingBlock, error)

func (sr *StoryRepository) Delete(any) error

func (sr *StoryRepository) WholeStory(int) (models.StartingBlock, []models.CommonBlock, error)

func (sr *StoryRepository) CreatedBlocks(int) ([]models.CommonBlock, error)

func (sr *StoryRepository) Create(any) (int, error)

func (sr *StoryRepository) CreatedFBView(int) models.DialoguesData

func (sr *StoryRepository) DeleteWholeStory(int) error

func (sr *StoryRepository) EditBView(int) models.DialoguesData

func (sr *StoryRepository) Update(id int, data any) error

func (sr *StoryRepository) Get(model string, id int) (any, error)
