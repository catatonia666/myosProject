package sqlstore

import (
	"dialogue/internal/models"
	"encoding/json"
)

type StoryRepository struct {
	store *Store
}

// Create creates a new instance of a passed model(starting_blocks or common_blocks) and returns ID of a new instance or an error.
func (sr *StoryRepository) Create(model any) (id int, err error) {
	err = sr.store.db.Create(model).Error
	if err != nil {
		return 0, err
	}

	err = sr.store.db.Model(model).Select("id").Order("id DESC").Last(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Get reads an instance of a model(starting_blocks or common_blocks) with ID  and returns it or an error.
func (sr *StoryRepository) Get(model string, id int) (data any, err error) {
	switch model {
	case "starting_blocks":
		{
			var firstBlock models.StartingBlock
			err = sr.store.db.Model(&models.StartingBlock{}).Where("id = ?", id).Find(&firstBlock).Error
			if err != nil {
				return models.StartingBlock{}, err
			}
			return firstBlock, nil
		}
	case "common_blocks":
		{
			var block models.CommonBlock
			err = sr.store.db.Model(&models.CommonBlock{}).Where("id = ?", id).Find(&block).Error
			if err != nil {
				return models.CommonBlock{}, err
			}
			return block, nil
		}
	default:
		{
			return nil, models.ErrNoRecord
		}
	}
}

// Update updates an instance of a model with a data passed.
func (sr *StoryRepository) Update(id int, data any) (err error) {
	err = sr.store.db.Model(data).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes an instance of a model passed.
func (sr *StoryRepository) Delete(model any) (err error) {
	err = sr.store.db.Unscoped().Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteWholeStory deletes the whole story with ID passed.
func (sr *StoryRepository) DeleteWholeStory(id int) error {
	err := sr.store.db.Unscoped().Where("id = ?", id).Delete(&models.StartingBlock{}).Error
	if err != nil {
		return err
	}

	err = sr.store.db.Unscoped().Where("story_id = ?", id).Delete(&models.CommonBlock{}).Error
	if err != nil {
		return err
	}
	return nil
}

// StoriesToDisplay gathers storiesAmount latest stories that user is able to see and displays it at the home page.
func (sr *StoryRepository) StoriesToDisplay(userID int, storiesAmount int) (storiesToDisplay []models.StartingBlock, err error) {
	err = sr.store.db.Model(&models.StartingBlock{}).Where("(privacy = false) OR (privacy = true AND user_id = ?)", userID).Limit(storiesAmount).Order("id desc").Scan(&storiesToDisplay).Error
	if err != nil {
		return nil, err
	}
	return storiesToDisplay, nil
}

// CreatedBlocks collects new blocks created along with the starting block of a new story. It is used only in CreateStory method.
func (sr *StoryRepository) CreatedBlocks(blocksAmount int) (retrievedBlocks []models.CommonBlock, err error) {
	err = sr.store.db.Model(&models.CommonBlock{}).Select("id").Limit(blocksAmount).Order("id desc").Scan(&retrievedBlocks).Error
	if err != nil {
		return nil, err
	}
	return retrievedBlocks, nil
}

// WholeStory collects all blocks related to the story with ID including the starting one.
func (sr *StoryRepository) WholeStory(storyID int) (story models.StartingBlock, blocks []models.CommonBlock, err error) {
	err = sr.store.db.Where("story_id = ?", storyID).Order("id").Find(&blocks).Error
	if err != nil {
		return models.StartingBlock{}, nil, err
	}

	err = sr.store.db.Where("id = ?", storyID).First(&story).Error
	if err != nil {
		return models.StartingBlock{}, nil, err
	}
	return story, blocks, nil
}

func (sr *StoryRepository) CreatedFBView(id int) (data models.DialoguesData) {
	var (
		firstBlock models.StartingBlock
		options    []map[int]string
	)
	sr.store.db.First(&firstBlock, id)
	json.Unmarshal(firstBlock.Options, &options)
	data.StartingBlock = firstBlock
	data.OptionsToBlocks = options
	data.RelatedToStoryBlocks, _ = sr.retrieveBlocks(id)
	return data
}

func (sr *StoryRepository) EditBView(id int) (data models.DialoguesData) {
	sr.store.db.First(&data.CommonBlock, id)
	var result []map[int]string
	json.Unmarshal(data.CommonBlock.Options, &result)
	data.OptionsToBlocks = result
	data.RelatedToStoryBlocks, _ = sr.retrieveBlocks(data.CommonBlock.StoryID)
	return data
}

// retrieveBlocks is a helper that converts blocks collected with WholeStory method into one struct.
func (sr *StoryRepository) retrieveBlocks(id int) (retrievedBlocks models.RelatedToStoryBlocks, err error) {
	firstBlock, blocks, err := sr.WholeStory(id)
	if err != nil {
		return models.RelatedToStoryBlocks{}, err
	}

	retrievedBlocks.StartingBlock = firstBlock
	retrievedBlocks.OtherBlocks = blocks
	return retrievedBlocks, nil
}
