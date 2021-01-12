package api

import (
	"github.com/gastrodon/scyther/storage"

	"context"
	"net/http"
	"strconv"
	"strings"
)

func splitPath(it rune) (ok bool) {
	ok = it == '/'
	return
}

func purportedTarget(path string) (target string, exists bool) {
	var parts []string = strings.FieldsFunc(path, splitPath)

	if exists = len(parts) >= 2; exists {
		target = parts[1]
	}

	return
}

func resolveTarget(path string) (id string, exists bool, err error) {
	if id, exists = purportedTarget(path); !exists {
		return
	}

	if uuidRegex.MatchString(id) {
		if _, exists, err = storage.ReadQueue(id); exists || err != nil {
			return
		}
	}

	id, exists, err = storage.ResolveNameId(id)
	return
}

func requestWithTarget(request *http.Request, target string) (modified *http.Request) {
	modified = request.WithContext(
		context.WithValue(
			request.Context(),
			keyQueue,
			target,
		),
	)

	return
}

func ResolveQueueTarget(request *http.Request) (modified *http.Request, ok bool, code int, RMap map[string]interface{}, err error) {
	var id string
	if id, ok, err = resolveTarget(request.URL.Path); err != nil || !ok {
		code = 404
		RMap = targetNotFound
		return
	}

	ok = true
	modified = requestWithTarget(request, id)
	return
}

func ResolveQueueIndex(request *http.Request) (modified *http.Request, ok bool, code int, RMap map[string]interface{}, err error) {
	var parts []string = strings.FieldsFunc(request.URL.Path, splitPath)
	var index int
	if index, err = strconv.Atoi(parts[len(parts)-1]); err != nil {
		err = nil
		code = 404
		RMap = noMessage
		return
	}

	ok = true
	modified = request.WithContext(
		context.WithValue(
			request.Context(),
			keyIndex,
			index,
		),
	)

	return
}

func ValidateLength(request *http.Request) (_ *http.Request, ok bool, code int, RMap map[string]interface{}, _ error) {
	var length int64 = request.ContentLength
	if length <= 0 {
		code = 411
		RMap = lengthRequired
		return
	}

	if length > storage.MESSAGE_MAX_SIZE {
		code = 413
		RMap = messageTooLong
		return
	}

	ok = true
	return
}
