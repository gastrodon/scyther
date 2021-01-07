package storage

import (
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"database/sql"
	"testing"
)

func Test_ReadQueues(test *testing.T) {
	test.Cleanup(Clear)

	var index, population int = 0, 10
	for index != population {
		_, _ = database.Exec(WRITE_QUEUE, uuid.New().String(), nil, nil)
		index++
	}

	var queues []types.QueueGet
	var count int
	var err error
	if queues, count, err = ReadQueues(); err != nil {
		test.Fatal(err)
	}

	if count != population {
		test.Fatalf("incorrect count: %d != %d", count, population)
	}

	if count != len(queues) {
		test.Fatalf("inaccurate count: %d != %d", count, len(queues))
	}
}

func Test_ReadQueues_empty(test *testing.T) {
	test.Cleanup(Clear)

	var population int = 0
	var queues []types.QueueGet
	var count int
	var err error
	if queues, count, err = ReadQueues(); err != nil {
		test.Fatal(err)
	}

	if count != population {
		test.Fatalf("incorrect count: %d != %d", count, population)
	}

	if count != len(queues) {
		test.Fatalf("inaccurate count: %d != %d", count, len(queues))
	}
}
func Test_ReadQueue(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	_, _ = database.Exec(WRITE_QUEUE, id, nil, nil)

	var queue types.QueueGet
	var exists bool
	var err error
	if queue, exists, err = ReadQueue(id); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s doesn't exist", id)
	}

	if queue.ID != id {
		test.Fatalf("id incorrect, %s != %s", queue.ID, id)
	}
}

func Test_ReadQueue_noSuchQueue(test *testing.T) {
	var exists bool
	var err error
	if _, exists, err = ReadQueue(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Fatalf("queue of random id exists")
	}
}

func Test_ReadQueue_err(test *testing.T) {
	test.Cleanup(reconnect)
	database, _ = sql.Open("mysql", "")

	var err error
	if _, _, err = ReadQueue(uuid.New().String()); err == nil {
		test.Fatal("no err was returned")
	}
}

func Test_ResolveNameId(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	var name string = "resolve_me"
	_, _ = database.Exec(WRITE_QUEUE, id, &name, nil)

	var resolvedId string
	var exists bool
	var err error
	if resolvedId, exists, err = ResolveNameId(name); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s doesn't exist", name)
	}

	if resolvedId != id {
		test.Fatalf("incorrect id, %s != %s", resolvedId, id)
	}
}

func Test_ResolveNameId_noSuchQueue(test *testing.T) {
	var exists bool
	var err error
	if _, exists, err = ResolveNameId(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Fatalf("queue of random id exists")
	}
}

func Test_WriteQueue(test *testing.T) {
	test.Cleanup(Clear)

	var id string
	var err error
	if id, err = WriteQueue(types.QueuePost{nil, nil, false}); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if err = database.QueryRow(TEST_READ_QUEUE_ID_EXISTS, id).Scan(&exists); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s doesn't exist", id)
	}
}

func Test_WriteQueue_nameConflict(test *testing.T) {
	test.Cleanup(Clear)

	var name string = "conflict_me"
	var err error
	if _, err = WriteQueue(types.QueuePost{&name, nil, false}); err != nil {
		test.Fatal(err)
	}

	if _, err = WriteQueue(types.QueuePost{&name, nil, false}); err == nil {
		test.Fatalf("No error for duplicate queue %s", name)
	}
}

func Test_WriteQueue_err(test *testing.T) {
	test.Cleanup(reconnect)
	database, _ = sql.Open("mysql", "")

	var err error
	if err = WriteMessage(uuid.New().String(), uuid.New().String()); err == nil {
		test.Fatal("no err was returned")
	}
}

func Test_DeleteQueue(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	_, _ = database.Exec(WRITE_QUEUE, id, nil, nil)

	var err error
	if err = DeleteQueue(id); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if database.QueryRow(TEST_READ_QUEUE_ID_EXISTS).Scan(&exists); exists {
		test.Fatalf("queue %s still exists", id)
	}
}

func Test_DeleteQueue_noSuchQueue(test *testing.T) {
	var err error
	if err = DeleteQueue(uuid.New().String()); err != nil {
		test.Fatal(err)
	}
}
