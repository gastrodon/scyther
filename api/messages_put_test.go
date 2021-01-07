package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"bytes"
	"net/http"
	"testing"
)

func Test_PutMessage(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	var data *bytes.Buffer = bytes.NewBufferString(message)
	var code int
	var err error
	if code, _, err = PutMessage(newRequestForQueue("PUT", "/queues/"+id, data, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 204, test)

	var queue types.QueueGet
	queue, _, _ = storage.ReadQueue(id)

	if queue.Size != 1 {
		test.Fatalf("queue %s is empty after messaging", id)
	}

	var messageFetched string
	var exists bool
	if messageFetched, exists, _ = storage.ReadHead(id); !exists {
		test.Fatalf("queue %s has no message on its head", id)
	}

	messageOk(message, messageFetched, test)
}

func Test_PutMessage_tooLong(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message, partial string = "foobar", "foo"
	var data *bytes.Buffer = bytes.NewBufferString(message)
	var size int64 = 3
	var request *http.Request = newRequestForQueue("PUT", "/queues/"+id, data, id)
	request.ContentLength = size

	var code int
	var err error
	if code, _, err = PutMessage(request); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 204, test)

	var messageFetched string
	messageFetched, _, _ = storage.ReadHead(id)

	messageOk(partial, messageFetched, test)
}

func Test_PutMessage_tooShort(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = uuid.New().String()
	var data *bytes.Buffer = bytes.NewBufferString(message)
	var size int64 = 90
	var request *http.Request = newRequestForQueue("PUT", "/queues/"+id, data, id)
	request.ContentLength = size

	var code int
	var err error
	if code, _, err = PutMessage(request); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 204, test)

	var messageFetched string
	messageFetched, _, _ = storage.ReadHead(id)

	messageOk(message, messageFetched, test)
}

func Test_PutMessage_pastCapacity(test *testing.T) {
	test.Cleanup(storage.Clear)

	var capacity int = 0
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, &capacity, false})

	var data *bytes.Buffer = bytes.NewBuffer([]byte("foobar"))
	var code int
	var err error
	if code, _, err = PutMessage(newRequestForQueue("PUT", "/queues/"+id, data, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 406, test)
}
