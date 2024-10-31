package stats

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smira/go-statsd"
)

var statsdClient *statsd.Client

func InitStatsDClient() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

    // Set the time field format to UTC
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	statsdClient= statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400),
	statsd.MetricPrefix("web."))
	log.Info().Msg("Logging started")
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

func TimeDataBaseQuery(query string, start time.Time, end time.Time){
	statsdClient.PrecisionTiming("db.query_time."+query, time.Duration(time.Since(start).Milliseconds()))
}

func TimeS3Call(query string, start time.Time, end time.Time){
	statsdClient.PrecisionTiming("s3.response_time."+query, time.Duration(time.Since(start).Milliseconds()))
}
