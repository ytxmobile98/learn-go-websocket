package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"example.com/websocket/examples/chat/chat"
)

const (
	defaultPort uint16 = 8080
)

var curDir = func() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}()

func main() {
	flag.Parse()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)

	fmt.Printf("Listening on http://127.0.0.1:%d\n", defaultPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		switch basePath := filepath.Base(r.URL.Path); basePath {
		case "index.html", "index.css", "index.js":
			http.ServeFile(w, r, filepath.Join(curDir, "static", basePath))
			return
		default:
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
	}
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, filepath.Join(curDir, "static", "index.html"))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	chat.Serve(w, r)
}
