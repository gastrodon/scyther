package main

import (
	"github.com/gastrodon/groudon"
	"github.com/gastrodon/scyther/api"
	"github.com/gastrodon/scyther/storage"

	"net/http"
	"os"
)

func main() {
	storage.Connect(os.Getenv("SCYTHER_CONNECTION"))

	groudon.RegisterHandler("GET", "^/queues/?$", api.GetQueues)
	groudon.RegisterHandler("POST", "^/queues/?$", api.CreateQueue)

	http.HandleFunc("/", groudon.Route)
	http.ListenAndServe(":8000", nil)
}
