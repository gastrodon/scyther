package api

import (
	"github.com/gastrodon/scyther/storage"

	"io"
	"net/http"
	"os"
	"testing"
)

func TestMain(main *testing.M) {
	storage.Connect(os.Getenv("SCYTHER_CONNECTION"))
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

func codeOk(code, want int, test *testing.T) {
	if code != want {
		test.Errorf("code %d != wanted %d", code, want)
	}
}
