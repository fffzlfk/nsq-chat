package config

import "time"

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
	MaxMessageSize    = 512
	PongWait          = 3 * time.Second
	PingPeriod        = PongWait * 9 / 10
	WriteWait         = 3 * time.Second
	ReadTimeout       = 5 * time.Second
	WriteTimeout      = 5 * time.Second
)
