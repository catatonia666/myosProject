package store

type Store interface {
	User() UserRepository
	Story() StoryRepository
}
