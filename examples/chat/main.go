package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	defaultPort uint16 = 8080
)

func main() {
	flag.Parse()

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/ws", serveWs)

	err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WebSocket endpoint"))
}
