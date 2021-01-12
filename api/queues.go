package api

import (
	"github.com/gastrodon/groudon/v2"
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"

	"net/http"
)

func GetQueues(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var queues []types.QueueGet
	var count int
	if queues, count, err = storage.ReadQueues(); err != nil {
		return
	}

	code = 200
	RMap = map[string]interface{}{
		"queues": queues,
		"count":  map[string]int{"queues": count},
	}

	return
}

func CreateQueue(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var data types.QueuePost
	var external error
	if err, external = groudon.SerializeBody(request.Body, &data); err != nil {
		return
	}

	if external != nil {
		code = 400
		RMap = badRequest
		return
	}

	var id string
	if id, err = storage.WriteQueue(data); err != nil {
		if err == storage.ErrNameOccupied {
			err = nil
			code = 409
			RMap = conflict
		}

		return
	}

	code = 200
	RMap = map[string]interface{}{"id": id}
	return
}

func GetQueue(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var id string = request.Context().Value(keyQueue).(string)
	var queue types.QueueGet
	if queue, _, err = storage.ReadQueue(id); err != nil {
		return
	}

	code = 200
	RMap = map[string]interface{}{"queue": queue}
	return
}

func DeleteQueue(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 204
	err = storage.DeleteQueue(request.Context().Value(keyQueue).(string))
	return
}
