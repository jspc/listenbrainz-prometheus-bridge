package main

import (
	"bytes"
	"encoding/csv"
	"io"

	_ "embed"
)

var (
	//go:embed geo/long_lat.csv
	longLatData []byte

	countries map[string]countryCapital
)

type countryCapital struct {
	lat, long string
}

func init() {
	countries = make(map[string]countryCapital)
	r := csv.NewReader(bytes.NewBuffer(longLatData))

	line := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		line++
		if line <= 1 {
			continue
		}

		countries[record[0]] = countryCapital{
			lat:  record[1],
			long: record[2],
		}
	}
}
