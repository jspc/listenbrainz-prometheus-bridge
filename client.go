package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type artistCount struct {
	Name  string `json:"artist_name"`
	Count int64  `json:"listen_count"`
}

func (a artistCount) Key() string {
	return a.Name
}

type trackCount struct {
	Artist string `json:"artist_name"`
	Album  string `json:"release_name"`
	Name   string `json:"track_name"`
	Count  int64  `json:"listen_count"`
}

func (t trackCount) Key() string {
	return fmt.Sprintf("%s%s%s", t.Artist, t.Album, t.Name)
}

type response struct {
	Payload struct {
		Artists []artistCount `json:"artists"`
		Tracks  []trackCount  `json:"recordings"`
	} `json:"payload"`
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	token string
	c     httpClient
}

func NewClient(token string) (c Client, err error) {
	c.token = token
	c.c = http.DefaultClient

	return
}

func (c Client) UserArtists(username string) (artists []artistCount, err error) {
	r, err := c.do(username, "artists")
	if err != nil {
		return
	}

	return r.Payload.Artists, nil
}

func (c Client) UserTracks(username string) (tracks []trackCount, err error) {
	r, err := c.do(username, "recordings")
	if err != nil {
		return
	}

	return r.Payload.Tracks, nil
}

func (c Client) do(username, model string) (payload *response, err error) {
	payload = new(response)

	request, err := http.NewRequest(http.MethodGet, deriveUrl(username, model), nil)
	if err != nil {
		return
	}

	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))

	r, err := c.c.Do(request)
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		err = errors.New("model " + model + " returned " + r.Status)

		return
	}

	err = json.NewDecoder(r.Body).Decode(payload)

	return
}

func deriveUrl(username, model string) string {
	return fmt.Sprintf("https://api.listenbrainz.org/1/stats/user/%s/%s", username, model)
}
