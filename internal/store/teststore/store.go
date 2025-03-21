package teststore

import (
	"dialogue/internal/models"
	"dialogue/internal/store"

	_ "gorm.io/driver/postgres"
)

type Store struct {
	userRepository  *UserRepository
	storyRepository *StoryRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*models.User),
	}

	return s.userRepository
}

func (s *Store) Story() store.StoryRepository {
	if s.userRepository != nil {
		return s.storyRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*models.User),
	}

	return s.storyRepository
}
