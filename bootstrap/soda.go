package bootstrap

import (
	"encoding/json"
	"github.com/SebastiaanKlippert/go-soda"
	"log"
)

const apiEndpoint string = "https://data.sfgov.org/resource/wwmu-gmzc"

// MovieLocation - Represents data returned by SODA API about SF film locations
type MovieLocation struct {
	Year     string `json:"release_year"`
	Title    string `json:"title"`
	Location string `json:"locations"`
}

// Locations - Get list of movies and their locations from SODA APIs
func Locations() ([]MovieLocation, error) {
	getreq := soda.NewGetRequest(apiEndpoint, "")
	getreq.Format = "json"
	getreq.Query.AddOrder("title", false)

	offsetreq, err := soda.NewOffsetGetRequest(getreq)
	if err != nil {
		return nil, err
	}
	locations := make([]MovieLocation, 0)
	offsetreq.Add(1)
	go func() {
		defer offsetreq.Done()

		for {
			resp, err := offsetreq.Next(2000)
			if err == soda.ErrDone {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			results := make([]MovieLocation, 0)
			err = json.NewDecoder(resp.Body).Decode(&results)
			resp.Body.Close()
			if err != nil {
				log.Println(err)
				continue
			}
			locations = append(locations, results...)
		}
	}()
	offsetreq.Wait()
	return locations, nil
}
