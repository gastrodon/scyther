package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"net/http"
	"testing"
)

func Test_ResolveQueueTarget(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = ResolveQueueTarget(newRequest("GET", "/queues/"+id, nil)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Fatal("middleware isn't ok")
	}

	var contextId string
	if contextId, ok = request.Context().Value(keyQueue).(string); !ok {
		test.Fatalf(
			"request context value %v is %v",
			keyQueue,
			request.Context().Value(keyQueue),
		)
	}

	if contextId != id {
		test.Fatalf("id incorrect, %s != %s", contextId, id)
	}
}

func Test_ResolveQueueTarget_named(test *testing.T) {
	test.Cleanup(storage.Clear)

	var name string = "resolve_me"
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{&name, nil, false})

	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = ResolveQueueTarget(newRequest("GET", "/queues/"+name, nil)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Fatal("middleware isn't ok")
	}

	var contextId string
	if contextId, ok = request.Context().Value(keyQueue).(string); !ok {
		test.Fatalf(
			"request context value %v is %v",
			keyQueue,
			request.Context().Value(keyQueue),
		)
	}

	if contextId != id {
		test.Fatalf("id incorrect, %s != %s", contextId, id)
	}

	var exists bool
	if _, exists, _ = storage.ReadQueue(contextId); !exists {
		test.Fatalf("queue %s doesn't exist", contextId)
	}
}
func Test_ResolveQueueTarget_subRoute(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})

	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = ResolveQueueTarget(newRequest("GET", "/queues/"+id+"/sub/route/", nil)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Fatal("middleware isn't ok")
	}

	var contextId string
	if contextId, ok = request.Context().Value(keyQueue).(string); !ok {
		test.Fatalf(
			"request context value %v is %v",
			keyQueue,
			request.Context().Value(keyQueue),
		)
	}

	if contextId != id {
		test.Fatalf("id incorrect, %s != %s", contextId, id)
	}
}

func Test_ResolveQueueTarget_namedSubRoute(test *testing.T) {
	test.Cleanup(storage.Clear)

	var name string = "resolve_me"
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{&name, nil, false})

	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = ResolveQueueTarget(newRequest("GET", "/queues/"+name+"/sub/route/", nil)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Fatal("middleware isn't ok")
	}

	var contextId string
	if contextId, ok = request.Context().Value(keyQueue).(string); !ok {
		test.Fatalf(
			"request context value %v is %v",
			keyQueue,
			request.Context().Value(keyQueue),
		)
	}

	if contextId != id {
		test.Fatalf("id incorrect, %s != %s", contextId, id)
	}

	var exists bool
	if _, exists, _ = storage.ReadQueue(contextId); !exists {
		test.Fatalf("queue %s doesn't exist", contextId)
	}
}

func Test_ResolveQueueTarget_namedAsUUID(test *testing.T) {
	test.Cleanup(storage.Clear)

	var name string = uuid.New().String()
	var id string
	id, _ = storage.WriteQueue(types.QueuePost{&name, nil, false})

	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = ResolveQueueTarget(newRequest("GET", "/queues/"+name, nil)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Fatal("middleware isn't ok")
	}

	var contextId string
	if contextId, ok = request.Context().Value(keyQueue).(string); !ok {
		test.Fatalf(
			"request context value %v is %v",
			keyQueue,
			request.Context().Value(keyQueue),
		)
	}

	if contextId != id {
		test.Fatalf("id incorrect, %s != %s", contextId, id)
	}

	var exists bool
	if _, exists, _ = storage.ReadQueue(contextId); !exists {
		test.Fatalf("queue %s doesn't exist", contextId)
	}
}

func Test_ResolveQueueTarget_nameConflictsId(test *testing.T) {
	test.Cleanup(storage.Clear)

	var id string
	id, _ = storage.WriteQueue(types.QueuePost{nil, nil, false})
	_, _ = storage.WriteQueue(types.QueuePost{&id, nil, false})

	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = ResolveQueueTarget(newRequest("GET", "/queues/"+id, nil)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Fatal("middleware isn't ok")
	}

	var contextId string
	if contextId, ok = request.Context().Value(keyQueue).(string); !ok {
		test.Fatalf(
			"request context value %v is %v",
			keyQueue,
			request.Context().Value(keyQueue),
		)
	}

	if contextId != id {
		test.Fatalf("id incorrect, %s != %s", contextId, id)
	}

	var exists bool
	if _, exists, _ = storage.ReadQueue(contextId); !exists {
		test.Fatalf("queue %s doesn't exist", contextId)
	}
}

func Test_ResolveQueueTarget_noSuchQueue(test *testing.T) {
	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = ResolveQueueTarget(newRequest("GET", "/queues/foobar", nil)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)

	if ok {
		test.Fatal("middleware is ok")
	}
}

func Test_ResolveQueueTarget_noTarget(test *testing.T) {
	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = ResolveQueueTarget(newRequest("GET", "/queues/foobar", nil)); err != nil {
		test.Fatal(err)
	}

	codeOk(code, 404, test)

	if ok {
		test.Fatal("middleware is ok")
	}
}
