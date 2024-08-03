package ws

import "time"

const (
	writeWait        = 10 * time.Second    // Time allowed to write a message to the peer.
	maxMessageSize   = 8192                // Maximum message size allowed from peer.
	pongWait         = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod       = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	closeGracePeriod = 10 * time.Second    // Time to wait before force close on connection.
)
