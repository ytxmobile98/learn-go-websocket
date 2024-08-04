package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/websocket/examples/filewatch/ws"
)

const homeHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(event) {
                    console.log('file updated');
                    data.textContent = event.data;
                }
            })();
        </script>
    </body>
</html>
`

var (
	addr = flag.String("addr", ":8080", "http service address")

	homeTempl = template.Must(template.New("").Parse(homeHTMLTemplate))
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("filename not specified")
	}

	filename := flag.Arg(0)

	http.HandleFunc("/", serveHome(filename))
	http.HandleFunc("/ws", serveWs(filename))

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func serveHome(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		header := w.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")
		bytes, lastMod, err := ws.ReadFileIfModified(filename, time.Time{})
		if err != nil {
			bytes = []byte(err.Error())
			lastMod = time.Unix(0, 0)
		}

		var v = struct {
			Host    string
			Data    string
			LastMod string
		}{
			Host:    r.Host,
			Data:    string(bytes),
			LastMod: strconv.FormatInt(lastMod.UnixNano(), 16),
		}

		homeTempl.Execute(w, &v)
	}
}

func serveWs(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.Serve(filename, w, r)
	}
}
