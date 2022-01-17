package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-01-01/2021-01-01"   | jq '.'
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-01-01/2021-01-01"   | jq '.data[]'
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-01-01/2021-01-01"   | jq '.countries[]'
	// Get and parse data:
	// link := buildLink()
	data := getData("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-01-17/2021-01-17")
	c := parseData(data)
	fmt.Println(c.Scale)
	fmt.Println("--- Full output:")
	fmt.Println(c)
}

type jsonBody struct {
	Scale struct {
		Deaths struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"deaths"`

		CasesConfirmed struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"casesConfirmed"`

		Stringency struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"stringency"`
	} `json:"scale"`
	Countries []string `json:"countries"`
}

type ParsedJson struct {
	Scale struct {
		Deaths struct {
			Min int
			Max int
		}
		CasesConfirmed struct {
			Min int
			Max int
		}
		Stringency struct {
			Min int
			Max int
		}
	}
	Countries []string
}

type Countries struct {
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-11-01/2021-11-12" | jq '.countries[]'
	Countries []string `json:"countries"`
}

func buildLink() string {
	var link string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range"
	tTime := time.Now()
	fmt.Printf("%s/%s-01-17/%s", link, tTime.Format("2006"), tTime.Format("2006-01-02"))
	return fmt.Sprintf("%s/%s-01-17/%s", link, tTime.Format("2006"), tTime.Format("2006-01-02"))
}

func getData(s string) []byte {
	response, err := http.Get(s)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func parseData(body []byte) ParsedJson {
	var jb jsonBody
	var pj ParsedJson
	err := json.Unmarshal(body, &jb)
	if err != nil {
		panic(err)
	}

	pj.Countries = jb.Countries
	// pj.Scale.Deaths = struct {
	// 	Min int
	// 	Max int
	// }(jb.Scale.Deaths)
	pj.Scale = struct {
		Deaths struct {
			Min int
			Max int
		}
		CasesConfirmed struct {
			Min int
			Max int
		}
		Stringency struct {
			Min int
			Max int
		}
	}(jb.Scale)

	return pj
}