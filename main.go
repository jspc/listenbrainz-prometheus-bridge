package main

import (
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	token = os.Getenv("LISTENBRAINZ_TOKEN")
	users = strings.Split(os.Getenv("LISTENBRAINZ_USERS"), ",")
)

func main() {
	c, err := NewClient(token)
	if err != nil {
		panic(err)
	}

	go Prometheus()

	err = doMetrics(c)
	if err != nil {
		panic(err)
	}

	for range time.Tick(time.Minute * 10) {
		err = doMetrics(c)
		if err != nil {
			panic(err)
		}
	}
}

func doMetrics(c Client) (err error) {
	artists := make(map[string]artistCount, 0)
	tracks := make(map[string]trackCount, 0)

	for _, user := range users {
		artistCounts, err := c.UserArtists(user)
		if err != nil {
			return err
		}

		for _, ac := range artistCounts {
			if _, ok := artists[ac.Key()]; !ok {
				artists[ac.Key()] = artistCount{
					Name: ac.Name,
				}
			}

			a := artists[ac.Key()]
			a.Count += ac.Count
			artists[ac.Key()] = a
		}

		trackCounts, err := c.UserTracks(user)
		if err != nil {
			return err
		}

		for _, tc := range trackCounts {
			if _, ok := tracks[tc.Key()]; !ok {
				tracks[tc.Key()] = trackCount{
					Artist: tc.Artist,
					Album:  tc.Album,
					Name:   tc.Name,
				}
			}

			t := tracks[tc.Key()]
			t.Count += tc.Count
			tracks[tc.Key()] = t
		}
	}

	for _, a := range artists {
		if a.Name == "" {
			continue
		}
		artistsGauge.With(prometheus.Labels{"artist": a.Name}).Set(float64(a.Count))
	}

	for _, t := range tracks {
		tracksGauge.With(prometheus.Labels{"artist": t.Artist, "track": t.Name, "album": t.Album}).Set(float64(t.Count))
	}

	return
}
