package main

import (
	"github.com/gastrodon/groudon"
	"github.com/gastrodon/scyther/api"
	"github.com/gastrodon/scyther/storage"
	"github.com/gastrodon/scyther/types"

	"net/http"
	"os"
)

const (
	ROUTE_QUEUES        = `^/queues/?$`
	ROUTE_QUEUE         = `^/queues/` + types.NAME_PATTERN + `/?$`
	ROUTE_QUEUE_HEAD    = `^/queues/` + types.NAME_PATTERN + `/head/?$`
	ROUTE_QUEUE_TAIL    = `^/queues/` + types.NAME_PATTERN + `/tail/?$`
	ROUTE_QUEUE_CONSUME = `^/queues/` + types.NAME_PATTERN + `/consume/[\d]+/?$`
	ROUTE_QUEUE_PEEK    = `^/queues/` + types.NAME_PATTERN + `/peek/[\d]+/?$`
	ROUTE_TARGETED      = `^/queues/` + types.NAME_PATTERN + `(/.*)?$`
)

var (
	connection string = os.Getenv("SCYTHER_CONNECTION")
)

func setup() {
	storage.Connect(connection)

	groudon.RegisterMiddlewareRoute([]string{"PUT"}, ROUTE_TARGETED, api.ValidateLength)
	groudon.RegisterMiddlewareRoute([]string{"GET"}, ROUTE_QUEUE_CONSUME, api.ResolveQueueIndex)
	groudon.RegisterMiddlewareRoute([]string{"GET"}, ROUTE_QUEUE_PEEK, api.ResolveQueueIndex)
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
}

func main() {
	setup()
	http.ListenAndServe(":8000", nil)
}
