package storage

import (
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type tableDescriptor struct {
	Name   string
	Fields string
	Schema string
}

var (
	database *sql.DB
)

func Connect(address string) {
	var err error
	if database, err = sql.Open("mysql", address); err != nil {
		panic(err)
	}

	if err = database.Ping(); err != nil {
		panic(err)
	}

	database.SetMaxOpenConns(150)
	create()
}

func ReadQueues() (queues []types.QueueGet, count int, err error) {
	count = countQueues()
	queues = make([]types.QueueGet, count)

	if count == 0 {
		return
	}

	var rows *sql.Rows
	if rows, err = database.Query(READ_MANY_QUEUES); err != nil || rows == nil {
		return
	}

	var index int
	for rows.Next() {
		var id string
		var name *string
		var capacity *int
		var size int
		if err = rows.Scan(&id, &name, &capacity, &size); err != nil {
			return
		}

		queues[index] = types.QueueGet{
			ID:        id,
			Name:      name,
			Capacity:  capacity,
			Size:      size,
			Ephemeral: false,
		}
		index++
	}

	return
}

func WriteQueue(queue types.QueuePost) (id string, err error) {
	id = uuid.New().String()
	_, err = database.Exec(WRITE_QUEUE, id, queue.Name, queue.Capacity)
	return
}
