package main

import (
	"github.com/gastrodon/groudon/v2"
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
	ROUTE_QUEUE_CONSUME = `^/queues/` + types.NAME_PATTERN + `/consume/(\d)+/?$`
	ROUTE_QUEUE_PEEK    = `^/queues/` + types.NAME_PATTERN + `/peek/(\d)+/?$`
	ROUTE_TARGETED      = `^/queues/` + types.NAME_PATTERN + `(/.*)?$`
)

var (
	connection string = os.Getenv("SCYTHER_CONNECTION")
)

func setup() {
	storage.Connect(connection)

	groudon.AddMiddleware("PUT", ROUTE_TARGETED, api.ValidateLength)
	groudon.AddMiddleware("GET", ROUTE_QUEUE_CONSUME, api.ResolveQueueIndex)
	groudon.AddMiddleware("GET", ROUTE_QUEUE_PEEK, api.ResolveQueueIndex)
	groudon.AddMiddleware("GET", ROUTE_TARGETED, api.ResolveQueueTarget)
	groudon.AddMiddleware("PUT", ROUTE_TARGETED, api.ResolveQueueTarget)
	groudon.AddMiddleware("DELETE", ROUTE_TARGETED, api.ResolveQueueTarget)

	groudon.AddHandler("GET", ROUTE_QUEUES, api.GetQueues)
	groudon.AddHandler("POST", ROUTE_QUEUES, api.CreateQueue)
	groudon.AddHandler("GET", ROUTE_QUEUE, api.GetQueue)
	groudon.AddHandler("DELETE", ROUTE_QUEUE, api.DeleteQueue)

	groudon.AddHandler("PUT", ROUTE_QUEUE, api.PutMessage)
	groudon.AddHandler("GET", ROUTE_QUEUE_HEAD, api.ConsumeHead)
	groudon.AddHandler("GET", ROUTE_QUEUE_TAIL, api.ConsumeTail)
	groudon.AddHandler("GET", ROUTE_QUEUE_CONSUME, api.ConsumeIndex)
	groudon.AddHandler("GET", ROUTE_QUEUE_PEEK, api.PeekIndex)

	http.HandleFunc("/", groudon.Route)
}

func main() {
	setup()
	http.ListenAndServe(":8000", nil)
}
