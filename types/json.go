package types

import (
	"github.com/gastrodon/groudon"

	"regexp"
)

const (
	UUID_PATTERN      = `[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`
	UUID_ONLY_PATTERN = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	NAME_PATTERN      = `[A-Za-z0-9-_]{1,255}`
	NAME_ONLY_PATTERN = `^[A-Za-z0-9-_]{1,255}$`
)

var (
	nameRegex *regexp.Regexp = regexp.MustCompile(NAME_ONLY_PATTERN)
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

func validOptionalName(it interface{}) (ok bool, err error) {
	if ok = it == nil; ok {
		return
	}

	var name string
	if name, ok = it.(string); !ok {
		return
	}

	ok = nameRegex.MatchString(name)
	return
}

func (QueuePost) Validators() (values map[string]func(interface{}) (bool, error)) {
	values = map[string]func(interface{}) (bool, error){
		"name":      validOptionalName,
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
