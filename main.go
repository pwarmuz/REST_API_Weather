package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Jeffail/gabs/v2"
)

type coordinate struct {
	Latitude  float64
	Longitude float64
}

type metaDate struct {
	Year  int
	Month int
	Day   int
}

const (
	urlTemplateFromCoordinates string = "https://www.metaweather.com/api//location/search/?lattlong=%f,%f"
	urlTemplateForTemperature  string = "https://www.metaweather.com/api/location/%d/%d/%d/%d"

	// cityUrls query is no longer valid...
	//cityUrls string = "https://public.opendatasoft.com/api/records/1.0/search/?dataset=1000-largest-us-cities-by-population-with-geographic-coordinates&facet=city&facet=state&sort=population&rows=100"
	// I've found this json which works
	urlForCities  string = "https://gist.githubusercontent.com/Miserlou/c5cd8364bf9b2420bb29/raw/2bf258763cdddd704f8ffd3ea9a3e81d25e2c6f6/cities.json"
	citiesLimited        = 2 // You can adjust this to a smaller # to limit the # of queries so the program doesn't run as long
)

/*
Challenge: Modify the go code below to calculate the average temperature in the 100 largest cities in the United States at the current time. Handle any errors if a city is missing temperature data and skip that city in the final calculation.
Requirements:
	1) calculate average temperature of the 100 largest cities
	2) use current date, cannot comply with "the current time" since URL Template for temperature requires date
		temperatures are rendered every 3 hours so current time will never match but the latest render is always in position 0 of the JSON
	3) handle errors if temperature is null
		* Does not say what to do with errors so I will print them. Null temperatures will not be included in the final calculation.
	4) skip city from calculation if temperature is null
*/

func main() {

	// Get city coordinates in order to get woeid
	cityData := doGetRequest(urlForCities)
	cityCoordinates := getCityCoordinates(cityData)

	// Requirement 2 date is current
	date := metaDate{
		Year:  time.Now().Year(),
		Month: int(time.Now().Month()),
		Day:   time.Now().Day(),
	}

	cityTemperatures := make([]float64, 0, citiesLimited)
	for i, cityCoords := range cityCoordinates {
		// Get woeid
		weatherCityData := doGetRequest(fmt.Sprintf(urlTemplateFromCoordinates, cityCoords.Latitude, cityCoords.Longitude))
		cityWoeid := getCityWoeid(weatherCityData)

		// Get temperatureData
		temperatureData := doGetRequest(getFormattedWeatherURL(cityWoeid, date))
		temp, err := getCurrentTemperatureForCoordinates(temperatureData)
		if err != nil {
			fmt.Println(err)
		} else {
			// Requirement 4, city is not appended to the array of cities with valid temperatures giving an accurate average
			fmt.Printf("Adding item %d with woeid %d with temperature of %.7f to list\n", i+1, cityWoeid, temp)
			cityTemperatures = append(cityTemperatures, temp)
		}
	}

	getAverage(cityTemperatures)
}

func doGetRequest(url string) []byte {
	res, err := http.Get(url)
	if err != nil { // handle http and io errors within function, no need to pass responsiblity when it can be handled here, this will just lead to more responsility passing
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}

// getCityCoordinates will parse data, verifying a valid latitude and longitude are present
func getCityCoordinates(cityData []byte) []coordinate {
	cityDataParsed, err := gabs.ParseJSON(cityData)
	// this error should be handled because we can't directly handle the errors within the third party api and should respond to it's errors
	if err != nil {
		panic(err)
	}

	cities := cityDataParsed.Children()
	cityCoordinates := make([]coordinate, 0, citiesLimited)

	for i, city := range cities {
		if i < citiesLimited {
			lat, ok := city.Path("latitude").Data().(float64)
			if !ok {
				// panic because latitudes are required to get a location
				log.Panicf("failed to retrieve latitude for entry %d", i)
			}
			long, ok := city.Path("longitude").Data().(float64)
			if !ok {
				// panic because longitutes are required to get a location
				log.Panicf("failed to retrieve longitute for entry %d", i)
			}
			cityCoordinates = append(cityCoordinates, coordinate{
				Latitude:  lat,
				Longitude: long,
			})
		} else {
			break
		}
	}
	return cityCoordinates
}

func getFormattedWeatherURL(cityWoeid int64, date metaDate) string {
	return fmt.Sprintf(urlTemplateForTemperature, cityWoeid, date.Year, date.Month, date.Day)
}

// getCityWoeid parses data to retrieve the woeid. There is an assumption the accurate woeid is located on the first entry
func getCityWoeid(cityData []byte) int64 {
	weatherCitiesParsed, err := gabs.ParseJSON(cityData)
	if err != nil {
		panic(err)
	}

	weatherCityWoeids, ok := weatherCitiesParsed.Path("0.woeid").Data().(float64)
	if !ok {
		// panic because woeid is requied for URL Template parsing
		log.Panic("failed to retrieve woeid")
	}
	return int64(weatherCityWoeids)
}

// getCurrentTemperatureForCoordinates returns a temperature or an error if null. There is an assumption the latest temperature is located on the first entry
func getCurrentTemperatureForCoordinates(temperatureData []byte) (float64, error) {
	weatherDataParsed, err := gabs.ParseJSON(temperatureData)
	if err != nil {
		panic(err)
	}
	value, ok := weatherDataParsed.Path("0.the_temp").Data().(float64)
	if !ok {
		// Requirement 3
		return 0.0, errors.New("failed to retrieve temperature, ommiting from temperature average")
	}
	return value, nil
}

func getAverage(temperatures []float64) float64 {
	var sum float64
	for _, temp := range temperatures {
		sum += temp
	}
	// Requirement 1
	average := sum / float64(len(temperatures))

	fmt.Println("The average temperature of the ", len(temperatures), " largest cities is ", average)
	return average
}
