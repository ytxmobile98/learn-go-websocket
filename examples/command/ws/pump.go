package ws

import (
	"bufio"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func pumpStdin(conn *websocket.Conn, w io.Writer) {
	defer conn.Close()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(handlePong(conn))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(conn *websocket.Conn, r io.Reader, done chan struct{}) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			conn.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	conn.Close()
}

func handlePong(conn *websocket.Conn) func(string) error {
	return func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	}
}
