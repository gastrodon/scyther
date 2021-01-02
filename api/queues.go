package api

import (
	"github.com/gastrodon/groudon"
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
		return
	}

	var id string
	if id, err = storage.WriteQueue(data); err != nil {
		return
	}

	code = 200
	RMap = map[string]interface{}{"id": id}
	return
}

func GetQueue(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	var id string = request.Context().Value(keyQueue).(string)
	var queue types.QueueGet
	var exists bool
	if queue, exists, err = storage.ReadQueue(id); err != nil {
		return
	}

	if !exists {
		code = 400
		RMap = targetNotFound
		return
	}

	code = 200
	RMap = map[string]interface{}{"queue": queue}
	return
}

func DeleteQueue(request *http.Request) (code int, RMap map[string]interface{}, err error) {
	code = 501
	RMap = map[string]interface{}{"error": "unimplemented"}
	return
}
