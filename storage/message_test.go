package storage

import (
	"github.com/google/uuid"

	"fmt"
	"strings"
	"testing"
)

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
	database.QueryRow(READ_MESSAGE_AT, id, 0).Scan(&messageFetched)
	messageOk(
		strings.TrimRight(messageFetched, string([]byte{0})),
		message,
		test,
	)
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
