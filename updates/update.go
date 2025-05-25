package updates

import "time"

type Update struct {
	PatternLength uint
	BytesReceived uint64
	BytesSent     uint64
	RunTime       time.Duration
}
