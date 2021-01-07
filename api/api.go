package api

import (
	"github.com/gastrodon/scyther/types"

	"regexp"
)

var (
	uuidRegex *regexp.Regexp = regexp.MustCompile(types.UUID_ONLY_PATTERN)
	nameRegex *regexp.Regexp = regexp.MustCompile(types.NAME_ONLY_PATTERN)

	keyQueue = key("queue")
	keyIndex = key("index")

	atCapacity     map[string]interface{} = map[string]interface{}{"error": "at_capacity"}
	badRequest     map[string]interface{} = map[string]interface{}{"error": "bad_request"}
	lengthRequired map[string]interface{} = map[string]interface{}{"error": "length_required"}
	messageTooLong map[string]interface{} = map[string]interface{}{"error": "message_too_long"}
	targetNotFound map[string]interface{} = map[string]interface{}{"error": "no_such_queue"}
	noMessage      map[string]interface{} = map[string]interface{}{"error": "no_message"}
	conflict       map[string]interface{} = map[string]interface{}{"error": "conflict"}
)

type key string
