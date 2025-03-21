package teststore

import "dialogue/internal/models"

type StoryRepository struct {
}

func (sr *StoryRepository) Latest(int) []models.FirstBlock {
	return nil
}

func (sr *StoryRepository) CreateFB(int, string, string, []string, bool) int {
	return 0
}

func (sr *StoryRepository) CreatedFBView(int) models.DialoguesData {
	return models.DialoguesData{}
}

func (sr *StoryRepository) EditFB(int, int, string, string, []string) {

}

func (sr *StoryRepository) DeleteFB(int) {

}

func (sr *StoryRepository) EditBView(int) models.DialoguesData {
	return models.DialoguesData{}
}
func (sr *StoryRepository) EditB(int, int, string, string, []string) {

}

func (sr *StoryRepository) DeleteB(int) {

}
