package main

import "dialogue/internal/store"

type StoryStruct struct {
	db store.Store
}

var stories *StoryStruct
