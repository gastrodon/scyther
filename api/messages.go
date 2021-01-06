package api

import (
	"github.com/gastrodon/scyther/storage"

	"bufio"
	"io"
	"net/http"
)

func PutMessage(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	// TODO do this in middleware
	var length int64 = request.ContentLength
	if length > storage.MESSAGE_MAX_SIZE {
		code = 413
		RMap = messageTooLong
		return
	}

	if length == -1 {
		code = 411
		RMap = lengthRequired
		return
	}

	var reader *bufio.Reader = bufio.NewReader(request.Body)

	var index int64
	var runes []rune = make([]rune, length)
	for index != length {
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
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func ConsumeTail(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func ConsumeIndex(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}

func PeekIndex(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}
