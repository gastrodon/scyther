package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"testing"
)

func Test_ConsumeHead(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = ConsumeHead(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)
	messageOk(message, RMap["message"].(string), test)
}

func Test_ConsumeHead_noMessages(test *testing.T) {
	var id string = uuid.New().String()
	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = ConsumeHead(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_ConsumeHead_consumes(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var err error
	if _, _, err = ConsumeHead(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	var code int
	if code, _, err = ConsumeHead(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
}

func Test_ConsumeTail(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	storage.WriteMessage(id, uuid.New().String())
	storage.WriteMessage(id, uuid.New().String())
	storage.WriteMessage(id, uuid.New().String())

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = ConsumeTail(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)
	messageOk(message, RMap["message"].(string), test)
}

func Test_ConsumeTail_consumes(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	storage.WriteMessage(id, uuid.New().String())
	storage.WriteMessage(id, uuid.New().String())
	storage.WriteMessage(id, uuid.New().String())

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var err error
	if _, _, err = ConsumeTail(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	var RMap map[string]interface{}
	if _, RMap, err = ConsumeTail(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	if message == RMap["message"].(string) {
		test.Fatalf("queue %s tail wasn't consumed", id)
	}
}

func Test_ConsumeTail_noMessages(test *testing.T) {
	var id string = uuid.New().String()
	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = ConsumeHead(newRequestForQueue("GET", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_ConsumeIndex(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = ConsumeIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/consume/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)
	messageOk(message, RMap["message"].(string), test)
}

func Test_ConsumeIndex_noMessages(test *testing.T) {
	var id string = uuid.New().String()
	var RMap map[string]interface{}
	var code int
	var err error
	if code, RMap, err = ConsumeIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/consume/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_ConsumeIndex_outOfBounds(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	storage.WriteMessage(id, uuid.New().String())

	var RMap map[string]interface{}
	var code int
	var err error
	if code, RMap, err = ConsumeIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/consume/10", nil, id, 10)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_ConsumeIndex_consumes(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var err error
	if _, _, err = ConsumeIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/consume/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	var RMap map[string]interface{}
	var code int
	if code, RMap, err = ConsumeIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/consume/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_PeekIndex(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	storage.WriteMessage(id, message)

	var code int
	var RMap map[string]interface{}
	var err error
	if code, RMap, err = PeekIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/peek/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 200, test)
	messageOk(message, RMap["message"].(string), test)
}

func Test_PeekIndex_noMessages(test *testing.T) {
	var id string = uuid.New().String()
	var RMap map[string]interface{}
	var code int
	var err error
	if code, RMap, err = PeekIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/peek/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_PeekIndex_outOfBounds(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	storage.WriteMessage(id, uuid.New().String())

	var RMap map[string]interface{}
	var code int
	var err error
	if code, RMap, err = PeekIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/peek/10", nil, id, 10)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)
	errorOk(RMap, noMessage, test)
}

func Test_PeekIndex_preserves(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})
	storage.WriteMessage(id, uuid.New().String())

	var first map[string]interface{}
	var err error
	if _, first, err = PeekIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/peek/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	var second map[string]interface{}
	if _, second, err = PeekIndex(newRequestForQueueIndex("GET", "/queues/"+id+"/peek/0", nil, id, 0)); err != nil {
		test.Fatal(err)
	}

	messageOk(first["message"].(string), second["message"].(string), test)
}
