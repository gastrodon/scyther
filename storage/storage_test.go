package storage

import (
	"github.com/google/uuid"

	"os"
	"testing"
)

var (
	connection string = os.Getenv("SCYTHER_CONNECTION_TEST_STORAGE")
)

func TestMain(main *testing.M) {
	Connect(connection)
	Clear()
	os.Exit(main.Run())
}

func Test_Connect_panic(test *testing.T) {
	test.Cleanup(reconnect)
	defer expectPanic(test)

	Connect("")
}

func messageOk(message, want string, test *testing.T) {
	if message != want {
		test.Fatalf("message incorrect, %s != %s", message, want)
	}
}

func seedQueue(id string, size int) {
	var index int
	for index != size {
		index++
		database.Exec(
			WRITE_MESSAGE,
			uuid.New().String(),
			id,
			uuid.New().String(),
		)
	}
}

func reconnect() {
	Connect(connection)
}

func expectPanic(test *testing.T) {
	if recover() == nil {
		test.Fatal("nothing was paniced")
	}
}
