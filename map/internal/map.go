package internal

import "fuxi/providee"

var Map *mmap

type mmap struct {
	providee.Providee
}

func NewMap() *mmap {
	Map := &mmap{}
	providee.NewProvidee(2, "map")
	return Map
}
