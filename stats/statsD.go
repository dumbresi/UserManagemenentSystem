package stats

import (
	"time"

	"github.com/smira/go-statsd"
)

var statsdClient *statsd.Client

func InitStatsDClient() {
	statsdClient= statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400),
	statsd.MetricPrefix("web."))
	
}

func CountAPICall(endpoint string) {
    statsdClient.Incr("api.calls."+endpoint, 1)
}

func TimeAPICall(endpoint string, start time.Time) {
    elapsed := time.Since(start).Milliseconds()
	statsdClient.Timing("api.response_time."+endpoint, int64(elapsed))
}
