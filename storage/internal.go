package storage

import (
	"database/sql"
)

func create() {
	if _, err := database.Exec(CREATE_TABLE_QUEUES); err != nil {
		panic(err)
	}

	if _, err := database.Exec(CREATE_TABLE_MESSAGES); err != nil {
		panic(err)
	}
}

func countQueues() (size int) {
	var row *sql.Row
	if row = database.QueryRow(COUNT_QUEUES); row != nil {
		row.Scan(&size)
	}

	return
}

func incrementSize(id string) (err error) {
	_, err = database.Exec(INCREMENT_QUEUE_SIZE, id)
	return
}

func decrementSize(id string) (err error) {
	_, err = database.Exec(DECREMENT_QUEUE_SIZE, id)
	return
}

func readAvailableCapacity(id string) (available int, capped bool, err error) {
	return
}

func queueHasSpace(id string) (available bool, err error) {
	var optionalAvailable *bool
	err = database.QueryRow(READ_QUEUE_HAS_CAPACITY, id).Scan(&optionalAvailable)
	available = optionalAvailable == nil || *optionalAvailable
	return
}
