package main

import "net/http"

func serveWs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WebSocket endpoint"))
}
