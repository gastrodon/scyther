package storage

import (
	"github.com/google/uuid"

	"fmt"
	"strings"
	"testing"
)

func Test_ReadIndex(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)

	var messageFetched string
	var exists bool
	var err error
	if messageFetched, exists, err = ReadIndex(id, 0, false); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s has no message on it at 0", id)
	}

	messageOk(messageFetched, message, test)
}

func Test_ReadIndex_noConsume(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)

	var messageFetched string
	var exists bool
	var err error
	if messageFetched, exists, err = ReadIndex(id, 0, false); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s has no message on it at 0", id)
	}

	messageOk(messageFetched, message, test)

	if _, exists, err = ReadIndex(id, 0, false); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s index 0 was consumed", id)
	}
}

func Test_ReadIndex_consume(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var messageFetched string
	var exists bool
	var err error
	if messageFetched, exists, err = ReadIndex(id, 0, true); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s has no message on it at 0", id)
	}

	messageOk(messageFetched, message, test)

	if _, exists, err = ReadIndex(id, 0, true); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Fatalf("queue %s at 0 wasn't consumed", id)
	}
}

func Test_ReadIndex_consumeUpdatesHead(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var first, second string = uuid.New().String(), uuid.New().String()
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, first)
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, second)
	database.Exec(INCREMENT_QUEUE_SIZE, id)
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var firstFetched, secondFetched string
	var err error
	if firstFetched, _, err = ReadIndex(id, 0, true); err != nil {
		test.Fatal(err)
	}

	if secondFetched, _, err = ReadIndex(id, 0, true); err != nil {
		test.Fatal(err)
	}

	messageOk(firstFetched, first, test)
	messageOk(secondFetched, second, test)
}

func Test_ReadIndex_consumeDecrements(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, uuid.New().String())
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var size int
	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 1, test)

	var err error
	if _, _, err = ReadIndex(id, 0, true); err != nil {
		test.Fatal(err)
	}

	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 0, test)
}

func Test_ReadIndex_noSuchQueue(test *testing.T) {
	var exists bool
	var err error
	if _, exists, err = ReadIndex(uuid.New().String(), 0, false); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Fatal("no such queue has a message on its head")
	}
}

func Test_ReadHead(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var messageFetched string
	var exists bool
	var err error
	if messageFetched, exists, err = ReadHead(id); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s has no message on its head", id)
	}

	messageOk(messageFetched, message, test)
}

func Test_ReadHead_headPreserve(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	seedQueue(id, 10)

	var messageFetched string
	var err error
	if messageFetched, _, err = ReadHead(id); err != nil {
		test.Fatal(err)
	}

	messageOk(messageFetched, message, test)
}

func Test_ReadHead_past(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	seedQueue(id, 10)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var messageFetched string
	var err error
	if messageFetched, _, err = ReadHead(id); err != nil {
		test.Fatal(err)
	}

	if messageFetched == message {
		test.Fatalf("queue %s head is new message %s", id, message)
	}
}

func Test_ReadHead_noSuchQueue(test *testing.T) {
	var exists bool
	var err error
	if _, exists, err = ReadHead(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Fatal("no such queue has a message on its head")
	}
}

func Test_ReadHead_decrements(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, uuid.New().String())
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var size int
	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 1, test)

	var err error
	if _, _, err = ReadHead(id); err != nil {
		test.Fatal(err)
	}

	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 0, test)
}

func Test_ReadTail(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	seedQueue(id, 10)

	var message string = "messagable"
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, message)
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var messageFetched string
	var exists bool
	var err error
	if messageFetched, exists, err = ReadTail(id); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Fatalf("queue %s has no message on its tail", id)
	}

	messageOk(messageFetched, message, test)
}

func Test_ReadTail_noSuchQueue(test *testing.T) {
	var exists bool
	var err error
	if _, exists, err = ReadTail(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Fatal("no such queue has a message on its head")
	}
}

func Test_ReadTail_decrements(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)
	database.Exec(WRITE_MESSAGE, uuid.New().String(), id, uuid.New().String())
	database.Exec(INCREMENT_QUEUE_SIZE, id)

	var size int
	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 1, test)

	var err error
	if _, _, err = ReadTail(id); err != nil {
		test.Fatal(err)
	}

	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 0, test)
}

func Test_WriteMessage(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var message string = "messagable"
	var err error
	if err = WriteMessage(id, message); err != nil {
		test.Fatal(err)
	}

	var messageFetched string
	database.QueryRow(READ_MESSAGE_DATA_AT, id, 0).Scan(&messageFetched)
	messageOk(
		strings.TrimRight(messageFetched, string([]byte{0})),
		message,
		test,
	)
}

func Test_WriteMessage_increments(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)
	WriteMessage(id, uuid.New().String())

	var size int
	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)
	sizeOk(size, 1, test)
}

func Test_WriteMessage_atCapacity(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	var size int = 10
	database.Exec(WRITE_QUEUE, id, nil, size)

	for size != 0 {
		WriteMessage(id, fmt.Sprintf("message %d", size))
		size--
	}

	var err error
	if err = WriteMessage(id, "overflow"); err != ErrAtCapacity {
		test.Fatal("overflowing doesn't return an err")
	}
}
