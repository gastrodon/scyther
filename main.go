package main

import (
	"github.com/gastrodon/groudon"
	"github.com/gastrodon/scyther/api"
	"github.com/gastrodon/scyther/storage"

	"net/http"
	"os"
)

const (
	pattern             = storage.QUEUE_NAME_PATTERN
	ROUTE_QUEUES        = `^/queues/?$`
	ROUTE_QUEUE         = `^/queues/` + storage.QUEUE_NAME_PATTERN + `/?$`
	ROUTE_QUEUE_HEAD    = `^/queues/` + storage.QUEUE_NAME_PATTERN + `/head/?$`
	ROUTE_QUEUE_TAIL    = `^/queues/` + storage.QUEUE_NAME_PATTERN + `/tail/?$`
	ROUTE_QUEUE_CONSUME = `^/queues/` + storage.QUEUE_NAME_PATTERN + `/consume/[\d]+/?$`
	ROUTE_QUEUE_PEEK    = `^/queues/` + storage.QUEUE_NAME_PATTERN + `/peek/[\d]+/?$`
	ROUTE_TARGETED      = `^/queues/` + storage.QUEUE_NAME_PATTERN + `(/.*)?$`
)

func main() {
	storage.Connect(os.Getenv("SCYTHER_CONNECTION"))

	groudon.RegisterMiddlewareRoute([]string{"GET", "PUT", "DELETE"}, ROUTE_TARGETED, api.ResolveQueueTarget)

	groudon.RegisterHandler("GET", ROUTE_QUEUES, api.GetQueues)
	groudon.RegisterHandler("POST", ROUTE_QUEUES, api.CreateQueue)
	groudon.RegisterHandler("GET", ROUTE_QUEUE, api.GetQueue)
	groudon.RegisterHandler("DELETE", ROUTE_QUEUE, api.DeleteQueue)

	groudon.RegisterHandler("PUT", ROUTE_QUEUE, api.PutMessage)
	groudon.RegisterHandler("GET", ROUTE_QUEUE_HEAD, api.ConsumeHead)
	groudon.RegisterHandler("GET", ROUTE_QUEUE_TAIL, api.ConsumeTail)
	groudon.RegisterHandler("GET", ROUTE_QUEUE_CONSUME, api.ConsumeIndex)
	groudon.RegisterHandler("GET", ROUTE_QUEUE_PEEK, api.PeekIndex)

	http.HandleFunc("/", groudon.Route)
	http.ListenAndServe(":8000", nil)
}
