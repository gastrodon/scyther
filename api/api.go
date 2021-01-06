package api

import (
	"regexp"
)

const (
	UUID_PATTERN = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
)

var (
	uuidRegex *regexp.Regexp = regexp.MustCompile(UUID_PATTERN)

	keyQueue = key("queue")

	atCapacity     map[string]interface{} = map[string]interface{}{"error": "at_capacity"}
	badRequest     map[string]interface{} = map[string]interface{}{"error": "bad_request"}
	lengthRequired map[string]interface{} = map[string]interface{}{"error": "length_required"}
	messageTooLong map[string]interface{} = map[string]interface{}{"error": "message_too_long"}
	targetNotFound map[string]interface{} = map[string]interface{}{"error": "no_such_queue"}
)

type key string
