package storage

import (
	"github.com/gastrodon/scyther/types"

	"testing"
)

func Test_ReadQueues(test *testing.T) {
	test.Cleanup(Clear)

	var index, population int = 0, 10
	for index != population {
		var queue types.QueuePost = types.QueuePost{
			Name:      nil,
			Capacity:  nil,
			Ephemeral: false,
		}

		_, _ = WriteQueue(queue)
		index++
	}

	var queues []types.QueueGet
	var count int
	var err error
	if queues, count, err = ReadQueues(); err != nil {
		test.Fatal(err)
	}

	if count != population {
		test.Fatalf("incorrect count: %d != %d", count, population)
	}

	if count != len(queues) {
		test.Fatalf("inaccurate count: %d != %d", count, len(queues))
	}
}
