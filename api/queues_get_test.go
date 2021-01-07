package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"testing"
)

func Test_GetQueues(test *testing.T) {
	test.Cleanup(storage.Clear)

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = GetQueues(newRequest("GET", "/queues", nil)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)

	var queueCount int = RMap["count"].(map[string]int)["queues"]
	if queueCount != 0 {
		test.Fatalf("queue count is %d", queueCount)
	}
}

func Test_GetQueues_populated(test *testing.T) {
	test.Cleanup(storage.Clear)

	var index, population int = 0, 10
	for index != population {
		index++
		var queue types.QueuePost = types.QueuePost{
			Name:      nil,
			Ephemeral: false,
			Capacity:  nil,
		}

		storage.WriteQueue(queue)
	}

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = GetQueues(newRequest("GET", "/queues", nil)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)

	var queueCount int = RMap["count"].(map[string]int)["queues"]
	if queueCount != population {
		test.Fatalf("queue count is %d, want population %d", queueCount, population)
	}

	var queues []types.QueueGet = RMap["queues"].([]types.QueueGet)
	if len(queues) != population {
		test.Fatalf("queue collection is size %d, want population %d", len(queues), population)
	}
}

func Test_GetQueues_ordered(test *testing.T) {
	test.Cleanup(storage.Clear)

	var index int = 10
	var written [10]string
	for index != 0 {
		index--
		var queue types.QueuePost = types.QueuePost{
			Name:      nil,
			Ephemeral: false,
			Capacity:  nil,
		}

		var id string
		id, _ = storage.WriteQueue(queue)
		written[index] = id
	}

	var RMap map[string]interface{}
	var err error
	if _, RMap, err = GetQueues(newRequest("GET", "/queues", nil)); err != nil {
		test.Fatal(err)
	}

	var queues []types.QueueGet = RMap["queues"].([]types.QueueGet)
	for index != len(queues) {
		if queues[index].ID != written[index] {
			test.Fatalf("queue mismatch at %d: written %s is not id %s", index, written, queues[index].ID)
		}

		index++
	}
}

func Test_GetQueue(test *testing.T) {
	test.Cleanup(storage.Clear)

	var name string = uuid.New().String()
	var capacity int = 420
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{&name, &capacity, false})

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = GetQueue(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)
	queueOk(RMap["queue"].(types.QueueGet), &name, &capacity, test)
}

func Test_GetQueue_err(test *testing.T) {
	test.Cleanup(reconnect)
	disconnect()

	var err error
	if _, _, err = GetQueue(newRequestForQueue("GET", "/queue", nil, "")); err == nil {
		test.Fatal("no err was retured")
	}
}
