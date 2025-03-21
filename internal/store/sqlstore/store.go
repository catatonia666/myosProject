package sqlstore

import (
	"dialogue/internal/store"

	"gorm.io/gorm"
)

type Store struct {
	db              *gorm.DB
	userRepository  *UserRepository
	storyRepository *StoryRepository
}

func New(db *gorm.DB) *Store {
	s := &Store{
		db: db,
	}

	s.storyRepository = &StoryRepository{store: s}
	s.userRepository = &UserRepository{store: s}

	return s
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Story() store.StoryRepository {
	if s.storyRepository != nil {
		return s.storyRepository
	}

	s.storyRepository = &StoryRepository{
		store: s,
	}

	return s.storyRepository
}
