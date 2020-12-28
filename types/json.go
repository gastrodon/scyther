package types

import (
	"github.com/gastrodon/groudon"
)

type QueueGet struct {
	ID        string  `json:"id"`
	Name      *string `json:"name"`
	Ephemeral bool    `json:"ephemeral"`
	Capacity  *int    `json:"capacity"`
	Size      int     `json:"size"`
}

type QueuePost struct {
	Name      *string `json:"name"`
	Ephemeral string  `json:"ephemeral"`
	Capacity  int     `json:"capacity"`
}

func (QueuePost) Validators() (values map[string]func(interface{}) (bool, error)) {
	values = map[string]func(interface{}) (bool, error){
		"name":      groudon.OptionalString,
		"ephemeral": groudon.OptionalBool,
		"capacity":  groudon.OptionalNumber,
	}

	return
}

func (QueuePost) Defaults() (values map[string]interface{}) {
	values = map[string]interface{}{
		"name":      "",
		"ephemeral": true,
		"capacity":  nil,
	}

	return
}
