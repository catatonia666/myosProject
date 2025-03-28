package services

import (
	"dialogue/internal/store"
)

type Service interface {
	User() UserService
	Story() StoryService
}

type ServiceStruct struct {
	db           store.Store
	userService  UserService
	storyService StoryService
}

func New(db store.Store) *ServiceStruct {
	return &ServiceStruct{
		db:           db,
		userService:  &UserStruct{db: db},
		storyService: &StoryStruct{db: db},
	}
}

func (s *ServiceStruct) User() UserService {
	return s.userService
}

func (s *ServiceStruct) Story() StoryService {
	return s.storyService
}
