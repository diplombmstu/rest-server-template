package sse_server

import (
    "fmt"
    "net/http"
    "github.com/golang/glog"
)

// extends Handler
type SseBroker struct {
    Notifier       chan SseEvent
    newClients     chan client
    closingClients chan client
    clients        map[string]client
}

func NewBroker() (broker *SseBroker) {
    broker = &SseBroker{
        Notifier:       make(chan SseEvent, 1),
        newClients:     make(chan client),
        closingClients: make(chan client),
        clients:        make(map[string]client),
    }

    go broker.listen()

    return
}

func (broker *SseBroker) ServeHTTP(rw http.ResponseWriter, req *http.Request, clientId string) {
    // Make sure that the writer supports flushing.
    flusher, ok := rw.(http.Flusher)

    if !ok {
        http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "text/event-stream")
    rw.Header().Set("Cache-Control", "no-cache")
    rw.Header().Set("Connection", "keep-alive")
    rw.Header().Set("Access-Control-Allow-Origin", "*")

    // Each connection registers its own message channel with the Broker's connections registry
    messageChan := make(chan []byte)

    client := client{clientId, messageChan}
    broker.newClients <- client

    // Remove this client from the map of connected clients when this handler exits.
    defer func() {
        broker.closingClients <- client
    }()

    // Listen to connection close and un-register messageChan
    notify := rw.(http.CloseNotifier).CloseNotify()

    go func() {
        <-notify
        broker.closingClients <- client
    }()

    for {
        // Write to the ResponseWriter
        // Server Sent Events compatible
        fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
        flusher.Flush()
    }
}

func (broker *SseBroker) listen() {
    for {
        select {
        case s := <-broker.newClients:
            broker.clients[s.Id] = s
            glog.Infof("Client added. %d registered clients", len(broker.clients))
            break
        case s := <-broker.closingClients:
            delete(broker.clients, s.Id)
            glog.Infof("Removed client. %d registered clients", len(broker.clients))
            break
        case event := <-broker.Notifier:
            routeEvent(broker, event)
            break
        }
    }
}

func routeEvent(broker *SseBroker, event SseEvent) {
    if event.Broadcast {
        for key := range broker.clients {
            client := broker.clients[key]
            client.DataChannel <- event.Data
        }

        return
    }

    if client, err := broker.clients[event.ReceiverId]; !err {
        client.DataChannel <- event.Data
    }
}
