package sse_server

type client struct {
    Id          string
    DataChannel chan []byte
}