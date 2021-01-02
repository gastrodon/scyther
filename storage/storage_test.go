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
