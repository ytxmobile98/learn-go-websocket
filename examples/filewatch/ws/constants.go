package ws

import "time"

const (
	writeWait  = 10 * time.Second    // time allowed to write the file to the client
	pongWait   = 60 * time.Second    // time allowed to read the next pong message from the client
	pingPeriod = (pongWait * 9) / 10 // send pings to client with this period; must be less than pongWait
	filePeriod = 10 * time.Second    // poll for file changes with this period
)
