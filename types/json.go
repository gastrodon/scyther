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
	Capacity  *int    `json:"capacity"`
	Ephemeral bool    `json:"ephemeral"`
}

func (QueuePost) Validators() (values map[string]func(interface{}) (bool, error)) {
	values = map[string]func(interface{}) (bool, error){
		"name":      groudon.OptionalString,
		"capacity":  groudon.OptionalNumber,
		"ephemeral": groudon.OptionalBool,
	}

	return
}

func (QueuePost) Defaults() (values map[string]interface{}) {
	values = map[string]interface{}{
		"name":      nil,
		"capacity":  nil,
		"ephemeral": true,
	}

	return
}
