package services

import (
	"dialogue/internal/models"
	"dialogue/internal/store"
	"encoding/json"
	"strconv"
	"strings"
)

type StoryStruct struct {
	db store.Store
}

type StoryService interface {
	Create(int, string, string, []string, bool) int
	Get(string, int) (any, error)
	Edit(string, int, int, string, string, []string) error
	DeleteOneBlock(id int) error
	DeleteWholeStory(id int) error
	DisplayStories(int) ([]models.StartingBlock, error)
	StartingBlockData(id int) (data models.DialoguesData)
	BlockData(id int) (data models.DialoguesData)
}

const storiesToDisplayAmount = 10

// Create creates
func (ss *StoryStruct) Create(userID int, storyTitle string, startingBlockContent string, startingBlockOptions []string, privacy bool) (storyID int) {
	var (
		newFirstBlock = &models.StartingBlock{
			UserID: userID,
		}
		blocksSlice  []models.CommonBlock //Store new created blocks related to fresh story.
		options      []map[int]string     //New first block option-childBlockID relationships.
		firstBlockID int                  //To store the ID of new story.
	)

	//Create empty first block and blocks that are related to options of first block.
	firstBlockID, err := ss.db.Story().Create(newFirstBlock)
	if err != nil {
		return
	}

	for range startingBlockOptions {
		var block = &models.CommonBlock{
			StoryID: firstBlockID,
			UserID:  userID,
		}
		ss.db.Story().Create(block)
	}

	//Get IDs of new created blocks, collect them, make option-block ID relationships and store it as json.
	blocksSlice, _ = ss.db.Story().CreatedBlocks(len(startingBlockOptions))
	reverseSlice(blocksSlice)
	for i, v := range blocksSlice {
		mapIDtitle := make(map[int]string)
		mapIDtitle[v.ID] = startingBlockOptions[i]
		options = append(options, mapIDtitle)
	}
	jsonData, _ := json.Marshal(options)

	//Gather all new data and update first block with it, then return resulting ID of the story.
	newFirstBlock = &models.StartingBlock{
		StoryTitle: storyTitle,
		Privacy:    privacy,
		Content:    startingBlockContent,
		Options:    jsonData,
	}
	// ss.db.Story().UpdateStory(firstBlockID, *newFirstBlock)
	ss.db.Story().Update(firstBlockID, newFirstBlock)
	return firstBlockID
}

func (ss *StoryStruct) Get(string, int) (any, error) {
	return nil, nil
}

func (ss *StoryStruct) Edit(model string, id int, userID int, storyTitle string, content string, newOptions []string) error {
	blockModelAny, err := ss.db.Story().Get(model, id)
	if err != nil {
		return err
	}
	switch model {
	case "starting_blocks":
		storyID := id
		var (
			startingBlock    models.StartingBlock
			retrievedOptions []map[int]string
		)
		startingBlock = blockModelAny.(models.StartingBlock)
		json.Unmarshal(startingBlock.Options, &retrievedOptions)
		newOptionsJSON := ss.recreateOptions(newOptions, retrievedOptions, storyID, userID)
		startingBlock = models.StartingBlock{
			StoryTitle: storyTitle,
			Content:    content,
			Options:    newOptionsJSON,
		}
		ss.db.Story().Update(storyID, startingBlock)

	case "common_blocks":
		var (
			commonBlock      models.CommonBlock
			retrievedOptions []map[int]string
		)
		commonBlock = blockModelAny.(models.CommonBlock)
		json.Unmarshal(commonBlock.Options, &retrievedOptions)
		newOptionsJSON := ss.recreateOptions(newOptions, retrievedOptions, commonBlock.StoryID, userID)
		commonBlock = models.CommonBlock{
			Content: content,
			Options: newOptionsJSON,
		}
		ss.db.Story().Update(id, commonBlock)
	}
	return nil
}

func (ss *StoryStruct) DeleteOneBlock(id int) (err error) {
	var block models.CommonBlock
	blockModelAny, err := ss.db.Story().Get("common_blocks", id)
	if err != nil {
		return err
	}
	block = blockModelAny.(models.CommonBlock)
	err = ss.deleteBlock(id, block.StoryID)
	if err != nil {
		return err
	}
	return nil
}

func (ss *StoryStruct) DeleteWholeStory(id int) (err error) {
	err = ss.db.Story().DeleteWholeStory(id)
	if err != nil {
		return err
	}
	return nil
}

func (ss *StoryStruct) DisplayStories(userID int) (storiesToDisplay []models.StartingBlock, err error) {
	storiesToDisplay, err = ss.db.Story().StoriesToDisplay(userID, storiesToDisplayAmount)
	if err != nil {
		return nil, err
	}
	return storiesToDisplay, nil
}

func (ss *StoryStruct) StartingBlockData(id int) (data models.DialoguesData) {
	var (
		firstBlock models.StartingBlock
		options    []map[int]string
	)
	firstBlockAny, _ := ss.db.Story().Get("starting_blocks", id)
	firstBlock = firstBlockAny.(models.StartingBlock)
	json.Unmarshal(firstBlock.Options, &options)
	data.StartingBlock = firstBlock
	data.OptionsToBlocks = options

	data.RelatedToStoryBlocks, _ = ss.db.Story().RetrieveBlocks(id)

	return data
}

func (ss *StoryStruct) BlockData(id int) (data models.DialoguesData) {
	var (
		block   models.CommonBlock
		options []map[int]string
	)
	blockAny, _ := ss.db.Story().Get("common_blocks", id)
	block = blockAny.(models.CommonBlock)
	json.Unmarshal(block.Options, &options)
	data.CommonBlock = block
	data.OptionsToBlocks = options

	data.RelatedToStoryBlocks, _ = ss.db.Story().RetrieveBlocks(block.StoryID)

	return data
}

// recreateOptions recreating options of the starting (first) block or other blocks of the story.
func (ss *StoryStruct) recreateOptions(blockOptions []string, retrievedOptions []map[int]string, id, userID int) []byte {
	for _, v := range blockOptions {
		command, newOption, _ := strings.Cut(v, " ")
		switch command {

		//add keyword adds a new option to the block.
		case "add":
			var block = &models.CommonBlock{
				StoryID: id,
				UserID:  userID,
			}
			ss.db.Story().Create(block)
			newOpt := make(map[int]string)
			newOpt[block.ID] = newOption
			retrievedOptions = append(retrievedOptions, newOpt)

		//addTo keyword adds an option that leads to an existing block.
		case "addTo":
			idString, text, _ := strings.Cut(newOption, " ")
			id, _ := strconv.Atoi(idString)
			newOpt := make(map[int]string)
			newOpt[id] = text
			retrievedOptions = append(retrievedOptions, newOpt)

		//change keyword changes an existing option and does not affect to what block it related to.
		case "change":
			idString, newOption, _ := strings.Cut(newOption, " ")
		lookingTroughSlice:
			for _, k := range retrievedOptions {
				id, _ := strconv.Atoi(idString)
				_, ok := k[id]
				if ok {
					k[id] = newOption
					break lookingTroughSlice
				}
			}

		//delete deletes block with ID and it's appearences in other blocks.
		case "delete":
			idString, _, _ := strings.Cut(newOption, " ")
			id, _ := strconv.Atoi(idString)
			var storyID models.CommonBlock
			result, err := ss.db.Story().Get("common_blocks", id)
			if err != nil {
				continue
			}
			storyID = result.(models.CommonBlock)
			ss.deleteBlock(id, storyID.StoryID)
			retrievedOptions = remove(retrievedOptions, id)
		default:
			continue
		}
	}
	jsonData, _ := json.Marshal(retrievedOptions)
	return jsonData
}

// deleteBlock deletes block with ID and all blocks related to it if they no longer have connections to other blocks.
func (ss *StoryStruct) deleteBlock(targetID, storyID int) (err error) {
	var allBlocks []models.CommonBlock
	_, allBlocks, err = ss.db.Story().WholeStory(storyID)
	if err != nil {
		return err
	}
	parentCount := make(map[int]int)
	for _, block := range allBlocks {
		var unmarshaledOpts []map[int]string
		json.Unmarshal(block.Options, &unmarshaledOpts)
		for _, v := range unmarshaledOpts {
			var id int
			for key := range v {
				id = key
				break
			}
			parentCount[id]++
		}
	}
	var cascadeDelete func(int)
	cascadeDelete = func(blockID int) {
		var block = models.CommonBlock{}
		// block, err := ss.db.Story().GetBlock(blockID)
		result, err := ss.db.Story().Get("common_blocks", blockID)
		if err != nil {
			return
		}
		block = result.(models.CommonBlock)
		ss.db.Story().Delete(&block)
		var unmarshaledOpts2 []map[int]string
		json.Unmarshal(block.Options, &unmarshaledOpts2)
		for _, childID := range unmarshaledOpts2 {
			var id int
			for key := range childID {
				id = key
				break
			}
			parentCount[id]--
			if parentCount[id] == 0 {
				cascadeDelete(id)
			}
		}
	}
	cascadeDelete(targetID)
	ss.clearOptions(targetID, storyID)
	return nil
}

// clearOptions searches for the block that was deleted to appear in other blocks' options and delete them.
func (ss *StoryStruct) clearOptions(id, storyID int) {
	firstBlock, relatedBlocks, _ := ss.db.Story().WholeStory(storyID)
	for _, b := range relatedBlocks {
		var unmarshaledOpts []map[int]string
		json.Unmarshal(b.Options, &unmarshaledOpts)

		newOpts := unmarshaledOpts
		for _, v := range unmarshaledOpts {
			for key := range v {
				if key == id {
					newOpts = remove(newOpts, id)
				}
			}
		}
		b.Options, _ = json.Marshal(newOpts)
		ss.db.Story().Update(b.ID, b)
	}

	firstBlockOptions := []map[int]string{}
	json.Unmarshal(firstBlock.Options, &firstBlockOptions)
	newFBOptions := firstBlockOptions
	for _, v := range firstBlockOptions {

		for key := range v {
			if key == id {
				newFBOptions = remove(newFBOptions, id)
			}
		}
	}
	firstBlock.Options, _ = json.Marshal(newFBOptions)
	ss.db.Story().Update(firstBlock.StoryID, firstBlock)

}

// remove removes one block from a map.
func remove(slice []map[int]string, id int) []map[int]string {
	var result int
Loop:
	for i, v := range slice {
		for key := range v {
			if key == id {
				result = i
				break Loop
			}
		}
	}
	return append(slice[:result], slice[result+1:]...)
}

// reverseSlice to reverse slice of blocks.
func reverseSlice(slice []models.CommonBlock) {
	n := len(slice)
	for i := 0; i < n/2; i++ {
		slice[i], slice[n-1-i] = slice[n-1-i], slice[i]
	}
}
