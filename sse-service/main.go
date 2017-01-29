package main

import (
    "net/http"
    "github.com/diplombmstu/rest-server-template/sse-service/sse_server"
    "time"
    "fmt"
    "log"
)

// just for debugging
func main() {
    broker := sse_server.NewBroker()

    go func() {
        for {
            time.Sleep(time.Second * 2)
            eventString := fmt.Sprintf("the time is %v", time.Now())
            log.Println("Receiving event")
            broker.Notifier <- *sse_server.NewBroadcastEvent([]byte(eventString))
        }
    }()

    http.ListenAndServe("localhost:3000", broker)
}
