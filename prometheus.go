package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	artistsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "artist_plays",
		Help: "Number of plays for a specific artist",
	}, []string{"artist"})

	tracksGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "track_plays",
		Help: "Number of plays for a specific tracks",
	}, []string{"track", "artist", "album"})
)

func Prometheus() {
	err := prometheus.DefaultRegisterer.Register(artistsGauge)
	if err != nil {
		panic(err)
	}

	err = prometheus.DefaultRegisterer.Register(tracksGauge)
	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              ":2112",
		ReadHeaderTimeout: time.Second,
	}

	panic(server.ListenAndServe())
}
