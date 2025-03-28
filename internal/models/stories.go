package models

import (
	"encoding/json"
	"time"
)

type StartingBlock struct {
	StoryTitle string `gorm:"type:text"`

	StoryID int `gorm:"column:id;primary_key"`
	UserID  int
	Privacy bool

	Content string          `gorm:"type:text"`
	Options json.RawMessage `gorm:"type:json;default:'{}'"`

	CreatedAt time.Time
}

type CommonBlock struct {
	ID      int `gorm:"primary_key"`
	StoryID int
	UserID  int

	Content string          `gorm:"type:text"`
	Options json.RawMessage `gorm:"type:json;default:'{}'"`

	CreatedAt time.Time
}

// DialoguesData is a collection of data that may be passed to templates.
type DialoguesData struct {
	StartingBlock   StartingBlock
	CommonBlock     CommonBlock
	OptionsToBlocks []map[int]string

	DialoguesToDisplay   []StartingBlock
	RelatedToStoryBlocks RelatedToStoryBlocks
}

type RelatedToStoryBlocks struct {
	StartingBlock StartingBlock
	OtherBlocks   []CommonBlock
}
