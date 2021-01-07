package main

import (
	"os"
	"testing"
)

func TestMain(main *testing.M) {
	connection = os.Getenv("SCYTHER_CONNECTION_TEST_API")
	os.Exit(main.Run())
}

func Test_setup(test *testing.T) {
	setup()
}
