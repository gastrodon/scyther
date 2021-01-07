package storage

import (
	"github.com/google/uuid"

	"os"
	"testing"
)

func TestMain(main *testing.M) {
	Connect(os.Getenv("SCYTHER_CONNECTION_TEST_STORAGE"))
	Clear()
	os.Exit(main.Run())
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
