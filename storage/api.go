package storage

import (
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	QUEUE_NAME_PATTERN = "[A-Za-z0-9-_]{1,255}"
)

var (
	database *sql.DB
)

type tableDescriptor struct {
	Name   string
	Fields string
	Schema string
}

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

func Clear() {
	database.Exec(CLEAR_TABLE_QUEUES)
	database.Exec(CLEAR_TABLE_MESSAGES)
}

func ReadQueues() (queues []types.QueueGet, count int, err error) {
	count = countQueues()
	queues = make([]types.QueueGet, count)

	if count == 0 {
		return
	}

	var rows *sql.Rows
	if rows, err = database.Query(READ_QUEUES); err != nil || rows == nil {
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

func ReadQueue(id string) (queue types.QueueGet, exists bool, err error) {
	var idRead string
	var name *string
	var capacity *int
	var size int
	if err = database.QueryRow(READ_QUEUE, id).Scan(&idRead, &name, &capacity, &size); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	queue = types.QueueGet{
		ID:        idRead,
		Name:      name,
		Capacity:  capacity,
		Size:      size,
		Ephemeral: false,
	}

	return
}

func ResolveNameId(name string) (id string, exists bool, err error) {
	if err = database.QueryRow(READ_QUEUE_ID, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	return
}

func WriteQueue(queue types.QueuePost) (id string, err error) {
	id = uuid.New().String()
	_, err = database.Exec(WRITE_QUEUE, id, queue.Name, queue.Capacity)
	return
}

func DeleteQueue(id string) (err error) {
	return
}

func ReadHead(id string, consume bool) (data []byte, exists bool, err error) {
	data, exists, err = ReadIndex(id, 0, consume)
	return
}

func ReadTail(id string, consume bool) (data []byte, exists bool, err error) {
	return
}

func ReadIndex(id string, index int, consume bool) (data []byte, exists bool, err error) {
	return
}

func dropMessage(id string) (err error) {
	return
}
