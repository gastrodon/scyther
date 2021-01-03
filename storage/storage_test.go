package storage

import (
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
