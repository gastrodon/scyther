package storage

import (
	"github.com/gastrodon/scyther/types"
	"github.com/google/uuid"

	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MESSAGE_MAX_SIZE = 255
)

var (
	database      *sql.DB
	ErrAtCapacity = errors.New("This message overflows the queue")
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
	if count = countQueues(); count == 0 {
		return
	}

	queues = make([]types.QueueGet, count)
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
	if _, err = database.Exec(DELETE_QUEUE, id); err == nil {
		_, err = database.Exec(DELETE_MESSAGES, id)
	}

	return
}

func ReadHead(queue string) (data string, exists bool, err error) {
	data, exists, err = ReadIndex(queue, 0, true)
	return
}

func ReadTail(queue string) (data string, exists bool, err error) {
	var id string
	if err = database.QueryRow(READ_TAIL, queue).Scan(&id, &data); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	_, err = database.Exec(DELETE_MESSAGE, id)
	return
}

func ReadIndex(queue string, index int, consume bool) (data string, exists bool, err error) {
	var id string
	if err = database.QueryRow(READ_MESSAGE_AT, queue, index).Scan(&id, &data); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	if consume {
		_, err = database.Exec(DELETE_MESSAGE, id)
	}

	return
}

func WriteMessage(queue string, message string) (err error) {
	var available bool
	if available, err = queueHasSpace(queue); err != nil {
		return
	}

	if !available {
		err = ErrAtCapacity
		return
	}

	if _, err = database.Exec(WRITE_MESSAGE, uuid.New().String(), queue, message); err == nil {
		err = incrementSize(queue)
	}

	return
}
