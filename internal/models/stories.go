package models

import (
	"encoding/json"
)

type FirstBlock struct {
	StoryTitle string `gorm:"type:text"`
	UserID     int
	ID         int `gorm:"primary_key"`
	Privacy    bool

	FirstBlockContent string          `gorm:"type:text"`
	FirstBlockOptions json.RawMessage `gorm:"type:json;default:'{}'"`
}

type Block struct {
	UserID  int
	StoryID int
	ID      int `gorm:"primary_key"`

	BlockContent string          `gorm:"type:text"`
	BlockOptions json.RawMessage `gorm:"type:json;default:'{}'"`
}

// DialoguesData is a collection of data that may be passed to templates.
type DialoguesData struct {
	FirstBlock      FirstBlock
	Block           Block
	OptionsToBlocks []map[int]string

	DialoguesToDisplay   []FirstBlock
	RelatedToStoryBlocks RelatedToStoryBlocks
}

type RelatedToStoryBlocks struct {
	FirstBlock  FirstBlock
	OtherBlocks []Block
}
