package client

import (
	"time"

	"github.com/mistsys/accord/protocol"
)

var (
	nanosInSecond = 1000000000
)

// calculate the latency from the metadata
// int64 in case the servers are out of sync
func latency(metadata *protocol.ReplyMetadata, curTime time.Time) (time.Duration, time.Duration) {
	respTimeNsecs := int64(metadata.ResponseTime.Seconds)*nanosInSecond + int64(metadata.ResponseTime.Nanos)
	reqTimeNsecs := int64(metadata.RequestTime.Seconds)*nanosInSecond + int64(metadata.RequestTime.Nanos)
	serverNsecs := respTimeNsecs - reqTimeNsecs
	totalDuration := curTime.Sub(time.Unix(0, reqTimeNsecs))
	return time.Duration(serverNsecs), totalDuration
}
