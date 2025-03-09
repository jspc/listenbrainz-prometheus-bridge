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

	geoGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "country_artists",
		Help: "Artists segregated by country",
	}, []string{"country", "lat", "long"})
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

	err = prometheus.DefaultRegisterer.Register(geoGauge)
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
