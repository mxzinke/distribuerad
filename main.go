package main

import (
	"distribuerad/events"
	"distribuerad/http"
	"flag"
	"fmt"
)

var (
	address = flag.String("address", "0.0.0.0", "To set the address which binds the server")
	port    = flag.Int("port", 3333, "To set the port, where the server starts")
)

func init() {
	flag.Parse()
}

func main() {
	bindAddr := fmt.Sprintf("%s:%d", *address, *port)
	events_http.StartHTTP(bindAddr, events.NewChannelStore())
}
