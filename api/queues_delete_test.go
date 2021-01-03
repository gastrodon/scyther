package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"

	"testing"
)

func Test_DeleteQueue(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})
	_ = storage.WriteMessage(id, "foobar")

	var code int
	var err error
	if code, _, err = DeleteQueue(newRequestForQueue("DELETE", "/queues/"+id, nil, id)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 204, test)

	var exists bool
	if _, exists, _ = storage.ReadQueue(id); exists {
		test.Fatalf("queue %s wasn't deleted", id)
	}

	if _, exists, _ = storage.ReadHead(id, false); exists {
		test.Fatalf("queue %s still has a message", id)
	}
}
