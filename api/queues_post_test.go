package api

import (
	"github.com/gastrodon/scyther/storage"

	"bytes"
	"net/http"
	"strings"
	"testing"
)

const (
	VALUED = iota + 1
	MISSING
	NIL
)

var (
	longName = strings.Repeat("0", 0xFF+1)
)

func Test_CreateQueue(test *testing.T) {
	test.Cleanup(storage.Clear)

	var names int = 3
	var capacities int = 3
	var ephemerals int = 2
	var queues []map[string]interface{} = make(
		[]map[string]interface{},
		names*capacities*ephemerals,
	)

	// queue permutations: 3 * 3 * 2
	// Name: 		&name,	<missing>, 	nil
	// Capacity: 	&cap,	<missing>, 	nil
	// Ephemeral: 	false,	<missing>
	var index int = 0
	for names != 0 {
		for capacities != 0 {
			for ephemerals != 0 {
				queues[index] = queuePermutation(names, capacities, ephemerals)
				ephemerals--
				index++
			}

			ephemerals = 2
			capacities--
		}

		capacities = 3
		names--
	}

	var queue map[string]interface{}
	for _, queue = range queues {
		var code int
		var RMap map[string]interface{}
		var err error
		if code, RMap, err = CreateQueue(newRequestMarshalled("POST", "/queues", queue)); err != nil {
			test.Fatal(err)
		}

		codeOk(code, 200, test)

		var id string = RMap["id"].(string)
		uuidOk(id, test)
	}
}

func Test_CreateQueue_badRequest(test *testing.T) {
	var requests []*http.Request = []*http.Request{
		newRequest("POST", "/queues", nil),
		newRequest("POST", "/queues", bytes.NewBuffer(make([]byte, 0))),
		newRequest("POST", "/queues", bytes.NewBuffer([]byte{255, 165, 69})),
		newRequestMarshalled("POST", "/queues", map[string]interface{}{"name": "foo bar with spaces"}),
		newRequestMarshalled("POST", "/queues", map[string]interface{}{"name": ""}),
		newRequestMarshalled("POST", "/queues", map[string]interface{}{"name": "dunno?"}),
		newRequestMarshalled("POST", "/queues", map[string]interface{}{"name": longName}),
		newRequestMarshalled("POST", "/queues", map[string]interface{}{"name": 42069}),
	}

	var index int = 0
	var request *http.Request
	for _, request = range requests {
		var code int
		var err error
		if code, _, err = CreateQueue(request); err != nil {
			test.Fatal(err)
		}

		if code != 400 {
			panic(index)
		}

		index++

		codeOk(code, 400, test)
	}
}

func Test_CreateQueue_err(test *testing.T) {
	test.Cleanup(reconnect)
	disconnect()

	var buffer *bytes.Buffer = bytes.NewBuffer([]byte("{}"))

	var code int
	var err error
	if code, _, err = CreateQueue(newRequest("POST", "/queues", buffer)); err == nil {
		test.Fatal(code)
		test.Fatal("no err was retured")
	}
}
