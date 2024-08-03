package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"

	"example.com/websocket/common"
	"example.com/websocket/examples/command/ws"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "http service address")

	curDir = common.GetCurDir()
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Must specify at least one argument")
	}

	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatal("LookPath:", err)
	}

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs(cmdPath, args))

	log.Fatal(http.ListenAndServe(*addr, nil))
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

func serveWs(cmdPath string, args []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.Serve(w, r, cmdPath, args)
	}
}
