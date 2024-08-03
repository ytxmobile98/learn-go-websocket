package ws

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Serve(w http.ResponseWriter, r *http.Request, cmdPath string, args []string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	outReader, outWriter, err := os.Pipe()
	if err != nil {
		internalError(conn, "stdout:", err)
		return
	}
	defer outReader.Close()
	defer outWriter.Close()

	inReader, inWriter, err := os.Pipe()
	if err != nil {
		internalError(conn, "stdin:", err)
		return
	}
	defer inReader.Close()
	defer inWriter.Close()

	proc, err := os.StartProcess(cmdPath, flag.Args(), &os.ProcAttr{
		Files: []*os.File{inReader, outWriter, outWriter},
	})
	if err != nil {
		internalError(conn, "StartProcess:", err)
		return
	}

	inReader.Close()
	outWriter.Close()

	stdoutDone := make(chan struct{})
	go pumpStdout(conn, outReader, stdoutDone)
	go ping(conn, stdoutDone)

	pumpStdin(conn, inWriter)

	// some commands will exit when stdin is closed
	inWriter.Close()

	// other commands need a bonk on the head
	if err := proc.Signal(os.Interrupt); err != nil {
		log.Println("Interrupt:", err)
	}

	select {
	case <-stdoutDone:
	case <-time.After(time.Second):
		// A bigger bonk on the head
		if err := proc.Signal(os.Kill); err != nil {
			log.Println("Kill:", err)
		}
		<-stdoutDone
	}

	if _, err := proc.Wait(); err != nil {
		log.Println("Wait:", err)
	}
}

func internalError(conn *websocket.Conn, message string, err error) {
	log.Println(message, err)
	conn.WriteMessage(websocket.TextMessage, []byte("Internal server error"))
}
