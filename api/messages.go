package api

import (
	"github.com/gastrodon/scyther/storage"

	"bufio"
	"io"
	"net/http"
)

func PutMessage(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var reader *bufio.Reader = bufio.NewReader(request.Body)
	var index int64
	var runes []rune = make([]rune, request.ContentLength)
	for index != request.ContentLength {
		var next rune
		if next, _, err = reader.ReadRune(); err != nil {
			if err == io.EOF {
				err = nil
				runes = runes[:index]
			}

			break
		}

		runes[index] = next
		index++
	}

	var id string = request.Context().Value(keyQueue).(string)
	if err = storage.WriteMessage(id, string(runes)); err != nil {
		if err == storage.ErrAtCapacity {
			code = 406
			RMap = atCapacity
			err = nil
		}

		return
	}

	code = 204
	return
}

func ConsumeHead(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var queue string = request.Context().Value(keyQueue).(string)
	var message string
	var exists bool
	message, exists, err = storage.ReadHead(queue)
	code, RMap = serveMessage(message, exists)
	return
}

func ConsumeTail(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var queue string = request.Context().Value(keyQueue).(string)
	var message string
	var exists bool
	message, exists, err = storage.ReadTail(queue)
	code, RMap = serveMessage(message, exists)
	return
}

func ConsumeIndex(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code, RMap, err = handleIndex(request, true)
	return
}

func PeekIndex(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code, RMap, err = handleIndex(request, false)
	return
}

func handleIndex(request *http.Request, consume bool) (code int, RMap map[string]interface{}, err error) {
	var queue string = request.Context().Value(keyQueue).(string)
	var index int = request.Context().Value(keyIndex).(int)
	var message string
	var exists bool
	message, exists, err = storage.ReadIndex(queue, index, consume)
	code, RMap = serveMessage(message, exists)
	return
}

func serveMessage(message string, exists bool) (code int, RMap map[string]interface{}) {
	if !exists {
		code = 404
		RMap = noMessage
		return
	}

	code = 200
	RMap = map[string]interface{}{"message": message}
	return
}
