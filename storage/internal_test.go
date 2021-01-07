package storage

import (
	"github.com/google/uuid"

	"database/sql"
	"testing"
)

func Test_create(test *testing.T) {
	test.Cleanup(reconnect)
	defer expectPanic(test)

	database, _ = sql.Open("mysql", "")
	create()
}

func Test_incrementSize(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var err error
	if err = incrementSize(id); err != nil {
		test.Fatal(err)
	}

	var size int
	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)

	if size != 1 {
		test.Fatalf("queue %s size wasn't increased", id)
	}
}

func Test_decrementSize(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	incrementSize(id)
	incrementSize(id)

	var err error
	if err = decrementSize(id); err != nil {
		test.Fatal(err)
	}

	var size int
	database.QueryRow(READ_QUEUE_SIZE, id).Scan(&size)

	if size != 1 {
		test.Fatalf("queue %s size wasn't increased", id)
	}
}

func Test_decrementSize_underflow(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var err error
	if err = decrementSize(id); err == nil {
		test.Fatalf("queue %s underflow returned no err", id)
	}
}

func Test_queueHasSpace(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, 100)

	var available bool
	var err error
	if available, err = queueHasSpace(id); err != nil {
		test.Fatal(err)
	}

	if !available {
		test.Fatalf("queue %s has no available capacity", id)
	}
}

func Test_queueHasSpace_nilAvailable(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, nil)

	var available bool
	var err error
	if available, err = queueHasSpace(id); err != nil {
		test.Fatal(err)
	}

	if !available {
		test.Fatalf("queue %s has no available capacity", id)
	}
}

func Test_queueHasSpace_filled(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, 3)

	incrementSize(id)
	incrementSize(id)
	incrementSize(id)

	var available bool
	var err error
	if available, err = queueHasSpace(id); err != nil {
		test.Fatal(err)
	}

	if available {
		test.Fatalf("queue %s has available capacity", id)
	}
}

func Test_queueHasSpace_zeroCapacity(test *testing.T) {
	test.Cleanup(Clear)

	var id string = uuid.New().String()
	database.Exec(WRITE_QUEUE, id, nil, 0)

	var available bool
	var err error
	if available, err = queueHasSpace(id); err != nil {
		test.Fatal(err)
	}

	if available {
		test.Fatalf("queue %s has available capacity", id)
	}
}
