package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"

	"bytes"
	"testing"
)

func Test_PutMessage(test *testing.T) {
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var message string = "foobar"
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
	if messageFetched, exists, _ = storage.ReadHead(id, false); !exists {
		test.Fatalf("queue %s has no message on its head", id)
	}

	messageOk(message, messageFetched, test)
}

func Test_PutMessage_tooLong(test *testing.T) {
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var size int = 512
	var data *bytes.Buffer = bytes.NewBuffer(nil)
	data.Grow(size)
	for size != 0 {
		size--
		data.WriteRune('0')
	}

	var code int
	var err error
	if code, _, err = PutMessage(newRequestForQueue("PUT", "/queues/"+id, data, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 413, test)
}

func Test_PutMessage_pastCapacity(test *testing.T) {
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
