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
