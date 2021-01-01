package api

import (
	"github.com/gastrodon/scyther/storage"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
)

const (
	UUID_PATTERN = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
)

var (
	uuidRegex *regexp.Regexp = regexp.MustCompile(UUID_PATTERN)
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

func newRequestMarshalled(method, path string, data interface{}) (request *http.Request) {
	var marshalled []byte
	var err error
	if marshalled, err = json.Marshal(data); err != nil {
		panic(err)
	}

	request = newRequest(method, path, bytes.NewBuffer(marshalled))
	return
}

func codeOk(code, want int, test *testing.T) {
	if code != want {
		test.Errorf("code %d != wanted %d", code, want)
	}
}

func uuidOk(uuid string, test *testing.T) {
	if !uuidRegex.MatchString(uuid) {
		test.Fatalf("%s doesn't match the pattern %s", uuid, UUID_PATTERN)
	}
}
