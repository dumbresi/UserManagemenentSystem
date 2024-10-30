package stats

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smira/go-statsd"
)

var statsdClient *statsd.Client

func InitStatsDClient() {
	statsdClient= statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400),
	statsd.MetricPrefix("web."))
	log.Info().Msg("StatsD client initlaized")
}

func CountAPICall(endpoint string) {
	log.Info().Msg("Count API increment")
    statsdClient.Incr("api.calls."+endpoint, 1)
}

func TimeAPICall(endpoint string, start time.Time) {
    elapsed := time.Since(start).Milliseconds()
	log.Info().Msg("Time API")
	statsdClient.Timing("api.response_time."+endpoint, int64(elapsed))
}
