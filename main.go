package main

import (
	"fmt"
	"net/http"
	"time"
)

type Notify struct {
	subscriber map[chan string]int
	message    chan string
	newClient  chan chan string
}

func main() {
	b := &Notify{
		make(map[chan string]int),
		make(chan string),
		make(chan (chan string)),
	}

	// Start processing events
	b.Start()

	http.Handle("/events/", b)

	go func() {
		for {

			// Create a little message to send to clients,
			// including the current time.
			b.message <- fmt.Sprintf("the time is %v", time.Now())

			<-time.After(1 * time.Second)

		}
	}()

	http.ListenAndServe(":9999", nil)
}

func (notify *Notify) Start() {

	go func() {
		for {
			select {

			case client := <-notify.newClient:

				notify.subscriber[client] = 1

			case message := <-notify.message:

				// push to all subscriber
				for s, _ := range notify.subscriber {
					s <- message
				}
			}
		}
	}()
}

// handler for http
func (n *Notify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageChan := make(chan string)

	addHeaders(w)
	n.newClient <- messageChan

	for {
		msg, open := <-messageChan
		if !open {
			break
		}

		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
	}
}

func addHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
}
