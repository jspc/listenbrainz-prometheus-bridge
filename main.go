package main

import (
	"log"
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
	log.Print("Loading metrics")

	artists := make(map[string]artistCount)
	tracks := make(map[string]trackCount)
	countryArtists := make(map[string]artistCountryCount)

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

		countryCounts, err := c.UserArtistMap(user)
		if err != nil {
			return err
		}

		for _, acc := range countryCounts {
			if _, ok := countryArtists[acc.Key()]; !ok {
				countryArtists[acc.Key()] = artistCountryCount{
					Country: acc.Country,
					lat:     acc.lat,
					long:    acc.long,
				}
			}

			ca := countryArtists[acc.Key()]
			ca.Count += acc.Count
			countryArtists[acc.Key()] = ca
		}
	}

	for _, a := range artists {
		artistsGauge.With(prometheus.Labels{"artist": a.Name}).Set(float64(a.Count))
	}

	for _, t := range tracks {
		tracksGauge.With(prometheus.Labels{"artist": t.Artist, "track": t.Name, "album": t.Album}).Set(float64(t.Count))
	}

	for _, c := range countryArtists {
		geoGauge.With(prometheus.Labels{"country": c.Country, "lat": c.lat, "long": c.long}).Set(float64(c.Count))
	}

	return
}
