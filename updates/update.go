package updates

import "time"

type Update struct {
	PatternLengthIn  uint
	PatternLengthOut uint
	BytesReceived    uint64
	BytesSent        uint64
	StartTime        time.Time
}
