package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	urlStr := u.String()
	log.Printf("Connecting to %s", urlStr)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(urlStr, nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				return
			}
			log.Printf("Received [%d]: %s", messageType, message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if err := writeMessage(conn, []byte(t.String())); err != nil {
				return
			}
		case <-interrupt:
			interruptConn(conn, done)
			return
		}
	}
}

func writeMessage(conn *websocket.Conn, message []byte) error {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Write:", err)
		return err
	}
	return nil
}

func interruptConn(conn *websocket.Conn, done chan struct{}) error {
	log.Println("Interrupt")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Write close:", err)
		return err
	}

	select {
	case <-done:
	case <-time.After(time.Second):
	}
	return nil
}
