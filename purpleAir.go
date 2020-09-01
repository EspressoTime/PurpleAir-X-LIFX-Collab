package main

import (
	// "fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Scale struct {
	pmlo  float64
	pmhi  float64
	aqilo float64
	aqihi float64
}

func (s Scale) convert(pm float64) float64 {
	pmdiff := s.pmhi - s.pmlo
	aqidiff := s.aqihi - s.aqilo
	pmoff := pm - s.pmlo
	return aqidiff/pmdiff*pmoff + s.aqilo
}

func newScale(pmlo, pmhi, aqilo, aqihi float64) Scale {
	return Scale{
		pmlo:  pmlo,
		pmhi:  pmhi,
		aqilo: aqilo,
		aqihi: aqihi,
	}
}

func pmToAqi(pm float64) float64 {
	scales := []Scale{
		newScale(0.0, 12.0, 0.0, 50.0),
		newScale(12.0, 35.4, 50.0, 100.0),
		newScale(35.4, 55.4, 100.0, 150.0),
		newScale(55.4, 150.4, 150.0, 200.0),
		newScale(150.4, 250.4, 200.0, 300.0),
		newScale(250.4, 500.4, 300.0, 500.0),
	}

	for _, scale := range scales {
		if pm > scale.pmlo && pm <= scale.pmhi {
			return scale.convert(pm)
		}
	}

	return 0
}

func lrapa(pm float64) float64 {
	return pm*0.5 - 0.66
}

type Result struct {
	PM25 string `json:"PM2_5Value"`
}

type PurpleJSON struct {
	// Name string `json:"name"`
	Results []Result `json:"results"`
}

func fetchPA() float64 {
	resp, fetchErr := http.Get("https://www.purpleair.com/json?key=C4NXX8GBKK2TXUR9&show=37703")
	if fetchErr != nil {
		log.Fatal(fetchErr)
	}
	defer resp.Body.Close()
	body, fetchErr := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	// mockJson := []byte(
	// 	`{
	// 		"ignored":"field",
	// 		"results":[
	// 			{"PM2_5Value":"110.27"},
	// 			{"PM2_5Value":"95.92"}
	// 		]
	// 	}`)

	pj := PurpleJSON{}
	err := json.Unmarshal(body, &pj)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(pj)
	if len(pj.Results) < 2 {
		log.Fatal("too few results")
	}

	r1, err := strconv.ParseFloat(pj.Results[0].PM25, 64)
	if err != nil {
		log.Fatal(err)
	}

	r2, err := strconv.ParseFloat(pj.Results[1].PM25, 64)
	if err != nil {
		log.Fatal(err)
	}

	pm := lrapa((r1 + r2) / 2)
	finalAQI := pmToAqi(pm)
	return finalAQI
}
