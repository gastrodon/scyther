package api

import (
	"github.com/gastrodon/scyther/storage"
	"github.com/google/uuid"

	"bytes"
	"net/http"
	"testing"
)

const (
	VALUED = iota + 1
	MISSING
	NIL
)

func postQueuePermutation(named, capped, ephemeral int) (queue map[string]interface{}) {
	queue = make(map[string]interface{})

	switch named {
	case VALUED:
		queue["name"] = uuid.New().String()
	case NIL:
		queue["name"] = nil
	}

	switch capped {
	case VALUED:
		queue["capacity"] = 512
	case NIL:
		queue["capacity"] = nil

	}

	switch ephemeral {
	case VALUED:
		queue["ephemeral"] = false
	}

	return
}

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
				queues[index] = postQueuePermutation(names, capacities, ephemerals)
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
		// TODO: https://github.com/gastrodon/groudon/issues/7
		// newRequest("POST", "/queues", nil),
		newRequest("POST", "/queues", bytes.NewBuffer(make([]byte, 0))),
		newRequest("POST", "/queues", bytes.NewBuffer([]byte{255, 165, 69})),
	}

	var request *http.Request
	for _, request = range requests {
		var code int
		var err error
		if code, _, err = CreateQueue(request); err != nil {
			test.Fatal(err)
		}

		codeOk(code, 400, test)
	}
}
