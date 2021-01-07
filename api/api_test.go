package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestMain(main *testing.M) {
	storage.Connect(os.Getenv("SCYTHER_CONNECTION_TEST_API"))
	storage.Clear()
	os.Exit(main.Run())
}

func newRequest(method, path string, data io.Reader) (request *http.Request) {
	var err error
	if request, err = http.NewRequest(method, "http://localhost"+path, data); err != nil {
		panic(err)
	}

	return
}

func newRequestMarshalled(method, path string, data interface{}) (request *http.Request) {
	var marshalled []byte
	var err error
	if marshalled, err = json.Marshal(data); err != nil {
		panic(err)
	}

	request = newRequest(method, path, bytes.NewBuffer(marshalled))
	return
}

func newRequestForQueue(method, path string, data io.Reader, target string) (request *http.Request) {
	request = newRequest(method, path, data).WithContext(
		context.WithValue(context.Background(), keyQueue, target),
	)

	return
}

func newRequestForQueueIndex(method, path string, data io.Reader, target string, index int) (request *http.Request) {
	request = newRequestForQueue(method, path, data, target)
	request = request.WithContext(
		context.WithValue(request.Context(), keyIndex, index),
	)

	return
}

func queuePermutation(named, capped, ephemeral int) (queue map[string]interface{}) {
	queue = make(map[string]interface{})

	switch named {
	case VALUED:
		queue["name"] = uuid.New().String()
	case NIL:
		queue["name"] = nil
	}

	switch capped {
	case VALUED:
		queue["capacity"] = 512
	case NIL:
		queue["capacity"] = nil

	}

	switch ephemeral {
	case VALUED:
		queue["ephemeral"] = false
	}

	return
}

func codeOk(code, want int, test *testing.T) {
	if code != want {
		test.Fatalf("code %d != wanted %d", code, want)
	}
}

func uuidOk(uuid string, test *testing.T) {
	if !uuidRegex.MatchString(uuid) {
		test.Fatalf("%s isn't valid UUIDv4", uuid)
	}
}

func queueOk(queue types.QueueGet, name *string, capacity *int, test *testing.T) {
	if name != nil && queue.Name == nil {
		test.Fatalf("queue %s isn't named", queue.ID)
	}

	if name != nil && *queue.Name != *name {
		test.Fatalf("queue %s is misnamed, %s != %s", queue.ID, *name, *queue.Name)
	}

	if capacity != nil && queue.Capacity == nil {
		test.Fatalf("queue %s isn't capped", queue.ID)
	}

	if capacity != nil && *queue.Capacity != *capacity {
		test.Fatalf("queue %s is miscapped, %d != %d", queue.ID, *capacity, *queue.Capacity)
	}
}

func messageOk(message, want string, test *testing.T) {
	if message != want {
		test.Fatalf("message incorrect, %s != %s", message, want)
	}
}
