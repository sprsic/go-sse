package main

import (
	"fmt"
	"time"
)

type Notify struct {
	// one subs
	subscriber map[chan string]struct{}
	message    chan string
	newClient  chan chan string
}

func main() {

	b := &Notify{
		make(map[chan string]struct{}),
		make(chan string),
		make(chan (chan string)),
	}

	// Start processing events
	b.Process()
	//send channel
	go func(message chan<- string) {
		for {
			time.Sleep(time.Second * 1)
			//send msg every second
			message <- fmt.Sprintf("%v", time.Now())
		}
	}(b.message)

	// keep main alive
	b.run()
}

func (notify *Notify) Process() {
	go func() {
		for {
			select {
			case client := <-notify.newClient: // chan string
				fmt.Println("new client registred")
				notify.subscriber[client] = struct{}{}
			case message := <-notify.message: // string
				// push to all subscriber
				fmt.Println("push message to sub")
				// key is the subscriber
				for k, _ := range notify.subscriber {
					k <- message
				}
			}
		}
	}()
}

// handler for http
func (n *Notify) run() {

	messageChan := make(chan string)
	// initial client
	n.newClient <- messageChan
	for {
		select {
		case msg := <-messageChan:
			fmt.Println("Current time: ", msg)
		}
	}
}
