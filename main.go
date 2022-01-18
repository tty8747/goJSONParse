package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	data := getData("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-01-16/2021-01-17")

	// Outputs list of Countries
	// fmt.Println(getListOfCoutries(data))
	// Generate list of dates
	// genListOfDates()
	// Get all data
	resultMap, err := getAllData(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(resultMap)
}

// Makes struct for json
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
	Data      struct {
	} `json:"2021-01-16"`
}

// Makes endpoint struct
type parsedJson struct {
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

// Makes map from json
func getAllData(body []byte) (m map[string]interface{}, err error) {
	err = json.Unmarshal(body, &m)
	return
}

// Gets raw data
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

// Gets list of countires
func getListOfCoutries(body []byte) []string {

	var jb jsonBody

	err := json.Unmarshal(body, &jb)
	if err != nil {
		panic(err)
	}
	return []string(jb.Countries)
}

// Gets list of dates
func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// Gets list of dates
func genListOfDates() (listOfDates []string) {
	end := time.Now()
	start, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%s-%s-%s", end.Format("2006"), "01", "01"), time.Local)
	if err != nil {
		panic(err)
	}

	for rd := rangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		// fmt.Println(date.Format("2006-01-02"))
		listOfDates = append(listOfDates, date.Format("2006-01-02"))
	}
	fmt.Println(listOfDates)
	return listOfDates
}

func makeLink() string {
	var link string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range"
	t := time.Now()
	fmt.Printf("%s/%s-01-17/%s", link, t.Format("2006"), t.Format("2006-01-02"))
	return fmt.Sprintf("%s/%s-01-17/%s", link, t.Format("2006"), t.Format("2006-01-02"))
}

func parseData(body []byte) parsedJson {
	var jb jsonBody
	var pj parsedJson
	err := json.Unmarshal(body, &jb)
	if err != nil {
		panic(err)
	}

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
