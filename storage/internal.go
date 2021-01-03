package storage

import (
	"database/sql"
	"sort"
	"time"
)

type capacityBuffer struct {
	Limit        int
	limitHard    int
	size         int
	capacities   map[string]int
	sizes        map[string]int
	interactions map[string]int64
}

func (buffer capacityBuffer) SetCapacity(queue string, capacity int) {
	buffer.capacities[queue] = capacity

	if buffer.size+1 > buffer.limitHard {
		buffer.interact(queue)
		go buffer.dropStale()
		return
	}

	go buffer.interact(queue)
}

func (buffer capacityBuffer) GetCapacity(queue string) (capacity int, exists bool, err error) {
	if capacity, exists = buffer.capacities[queue]; exists {
		go buffer.interact(queue)
	}

	return
}

func (buffer capacityBuffer) interact(queue string) {
	buffer.interactions[queue] = time.Now().Unix()
}

func (buffer capacityBuffer) dropStale() {
	var size int = len(buffer.interactions)

	var stamps []int64 = make([]int64, size)
	var index int = 0

	var interactionsFlip map[int64]string = make(map[int64]string, size)
	var queue string
	var stamp int64
	for queue, stamp = range buffer.interactions {
		interactionsFlip[stamp] = queue
		stamps[index] = stamp
		index++
	}

	sort.Slice(stamps, func(it, next int) bool { return stamps[it] < stamps[next] })
	for _, stamp = range stamps[buffer.Limit:] {
		queue = interactionsFlip[stamp]
		delete(buffer.capacities, queue)
		delete(buffer.interactions, queue)
	}
}

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
