package sqlstore

import (
	"dialogue/internal/models"
	"encoding/json"
	"strconv"
	"strings"
)

type StoryRepository struct {
	store *Store
}

// CreateFB inserts starting block (FB - first block) into the database.
func (sr *StoryRepository) CreateFB(userid int, firstBlockTitle, firstBlockContent string, firstBlockOptions []string, privacy bool) (id int) {
	var (
		newFirstBlock models.FirstBlock //Variable to store data related to new first block of a story.
		blocksSlice   []models.Block    //Store new created blocks related to fresh story.
		options       []map[int]string  //New first block option-childBlockID relationships.
		firstBlockID  models.FirstBlock //To store the ID of new story.
	)

	//Create empty first block and blocks that are related to options of first block.
	sr.store.db.Create(&newFirstBlock)
	for range firstBlockOptions {
		var block models.Block = models.Block{
			StoryID: newFirstBlock.ID,
			UserID:  userid,
		}
		sr.store.db.Create(&block)
	}

	//Get IDs of new created blocks, collect them, make option-block ID relationships and store it as json.
	sr.store.db.Model(&models.Block{}).Select("ID").Limit(len(firstBlockOptions)).Order("id desc").Scan(&blocksSlice)
	reverseSlice(blocksSlice)
	for i, v := range blocksSlice {
		mapIDtitle := make(map[int]string)
		mapIDtitle[v.ID] = firstBlockOptions[i]
		options = append(options, mapIDtitle)
	}
	jsonData, _ := json.Marshal(options)

	//Gather all new data and update first block with it, then return resulting ID of the story.
	newFirstBlock = models.FirstBlock{
		StoryTitle:        firstBlockTitle,
		UserID:            userid,
		Privacy:           privacy,
		FirstBlockContent: firstBlockContent,
		FirstBlockOptions: jsonData,
	}
	sr.store.db.Select("ID").Last(&firstBlockID)
	sr.store.db.Model(&models.FirstBlock{}).Where("id = ?", firstBlockID.ID).Updates(&newFirstBlock)
	return firstBlockID.ID
}

// CreatedFBView gets the nessessary data related to fresh created story and pass it to render the view.
func (sr *StoryRepository) CreatedFBView(id int) (data models.DialoguesData) {
	var (
		firstBlock models.FirstBlock //Where data related to first block of fresh story be collected.
		options    []map[int]string  //Where unmarshalled options-id relationships stored.
	)

	//Get the fresh created story.
	sr.store.db.First(&firstBlock, id)

	//Unmarshall option-ID from json.
	json.Unmarshal(firstBlock.FirstBlockOptions, &options)

	//Gather all new data and pass it to render the view.
	data.FirstBlock = firstBlock
	data.OptionsToBlocks = options
	data.RelatedToStoryBlocks = sr.RetrieveBlocks(id)
	return data
}

// EditFB updates info about the first block user editing.
func (sr *StoryRepository) EditFB(id, userID int, blockTitle, blockContent string, blockOptions []string) {

	//Get all existing information about the first block that is about to be edited.
	var (
		editingFB        models.FirstBlock
		retrievedOptions []map[int]string
	)
	sr.store.db.Model(&models.FirstBlock{}).Where("id = ?", id).Find(&editingFB)
	idProviding := editingFB.ID
	json.Unmarshal(editingFB.FirstBlockOptions, &retrievedOptions)

	//Gather all new options for the block and update info.
	result := sr.recreateOptions(blockOptions, retrievedOptions, idProviding, userID)
	editingFB = models.FirstBlock{
		StoryTitle:        blockTitle,
		FirstBlockContent: blockContent,
		FirstBlockOptions: result,
	}
	sr.store.db.Model(&models.FirstBlock{}).Where("id = ?", id).Updates(&editingFB)
}

// DeleteFB deletes the whole story with ID.
func (sr *StoryRepository) DeleteFB(id int) {
	sr.store.db.Unscoped().Where("id = ?", id).Delete(&models.FirstBlock{})
	sr.store.db.Unscoped().Where("story_id = ?", id).Delete(&models.Block{})
}

// EditBView gets the data related to the block of the story and pass it to render.
func (sr *StoryRepository) EditBView(id int) (data models.DialoguesData) {

	//Get existing data about block that about to be edited.
	sr.store.db.First(&data.Block, id)
	var result []map[int]string
	json.Unmarshal(data.Block.BlockOptions, &result)
	data.OptionsToBlocks = result
	data.RelatedToStoryBlocks = sr.RetrieveBlocks(data.Block.StoryID)
	return data
}

// EditB update info about the block user editing.
func (sr *StoryRepository) EditB(id, userID int, blockTitle, blockContent string, blockOptions []string) {

	//Get all existing information about the block that is about to be edited.
	var (
		editingBlock     models.Block
		retrievedOptions []map[int]string
	)
	sr.store.db.Model(&models.Block{}).Where("id = ?", id).Find(&editingBlock)
	idProviding := editingBlock.StoryID
	json.Unmarshal(editingBlock.BlockOptions, &retrievedOptions)

	//Gather all new options for the block and update info.
	result := sr.recreateOptions(blockOptions, retrievedOptions, idProviding, userID)
	editingBlock = models.Block{
		BlockContent: blockContent,
		BlockOptions: result,
	}
	sr.store.db.Model(&models.Block{}).Where("id = ?", id).Updates(&editingBlock)
}

// DeleteB deletes block and it's appearances in other blocks with provided ID.
func (sr *StoryRepository) DeleteB(id int) {
	var block models.Block
	sr.store.db.Where("id = ?", id).First(&block)
	sr.deleteBlock(id, block.StoryID)
}

// RetrieveBlocks gets all blocks, including the starting one, that are parts of a story with ID.
func (sr *StoryRepository) RetrieveBlocks(id int) (retrievedBlocks models.RelatedToStoryBlocks) {
	sr.store.db.Model(&models.FirstBlock{}).Where("id = ?", id).First(&retrievedBlocks.FirstBlock)
	sr.store.db.Model(&models.Block{}).Where("story_id = ?", id).Order("id").Scan(&retrievedBlocks.OtherBlocks)
	return retrievedBlocks
}

// Latest gathers 10 latest stories that user is able to see and displays it at the home page.
func (sr *StoryRepository) Latest(userID int) (storiesToDisplay []models.FirstBlock) {
	sr.store.db.Table("first_blocks").Where("(privacy = false) OR (privacy = true AND user_id = ?)", userID).Limit(10).Order("id desc").Find(&storiesToDisplay)
	return storiesToDisplay
}

// reverseSlice to reverse slice of blocks.
func reverseSlice(slice []models.Block) {
	n := len(slice)
	for i := 0; i < n/2; i++ {
		slice[i], slice[n-1-i] = slice[n-1-i], slice[i]
	}
}

// recreateOptions recreating options of the starting (first) block or other blocks of the story.
func (sr *StoryRepository) recreateOptions(blockOptions []string, retrievedOptions []map[int]string, id, userID int) []byte {
	for _, v := range blockOptions {
		command, newOption, _ := strings.Cut(v, " ")
		switch command {

		//add keyword adds a new option to the block.
		case "add":
			var block models.Block = models.Block{
				StoryID: id,
				UserID:  userID,
			}
			sr.store.db.Create(&block)
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
			var storyID models.Block
			err := sr.store.db.Where("id = ?", id).Find(&storyID).Error
			if err != nil {
				continue
			}
			sr.deleteBlock(id, storyID.StoryID)
			retrievedOptions = remove(retrievedOptions, id)
		default:
			continue
		}
	}
	jsonData, _ := json.Marshal(retrievedOptions)
	return jsonData
}

// deleteBlock deletes block with ID and all blocks related to it if they no longer have connections to other blocks.
func (sr *StoryRepository) deleteBlock(targetID, storyID int) {
	var allBlocks []models.Block
	sr.store.db.Where("story_id = ?", storyID).Find(&allBlocks)
	parentCount := make(map[int]int)
	for _, block := range allBlocks {
		var unmarshaledOpts []map[int]string
		json.Unmarshal(block.BlockOptions, &unmarshaledOpts)
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
		var block models.Block
		if err := sr.store.db.First(&block, blockID).Error; err != nil {
			return
		}
		sr.store.db.Unscoped().Delete(&block)
		var unmarshaledOpts2 []map[int]string
		json.Unmarshal(block.BlockOptions, &unmarshaledOpts2)
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
	sr.clearOptions(targetID, storyID)
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

// clearOptions searches for the block that was deleted to appear in other blocks' options and delete them.
func (sr *StoryRepository) clearOptions(id, storyID int) {
	var relatedBlocks []models.Block
	sr.store.db.Model(&models.Block{}).Where("story_id = ?", storyID).Find(&relatedBlocks)
	for _, b := range relatedBlocks {
		var unmarshaledOpts []map[int]string
		json.Unmarshal(b.BlockOptions, &unmarshaledOpts)
		newOpts := unmarshaledOpts
		for _, v := range unmarshaledOpts {
			for key := range v {
				if key == id {
					newOpts = remove(newOpts, id)
				}
			}
		}
		sr.store.db.Model(&models.Block{}).Where("id = ?", b.ID).Update("block_options", &newOpts)
	}
}
