package main

import "testing"

func TestGetTemperatureFromCoords(t *testing.T) {
	temperatureRange := []struct {
		testData []byte
		wantTemp float64
	}{
		{
			testData: []byte(`[
				{
					"id":429009,
					"the_temp":27.67,
					"wind_speed":9.2608902
			}]`),
			wantTemp: 27.67,
		},
		{
			testData: []byte(`[
				{
					"id":429009,
					"the_temp":null,
					"wind_speed":9.2608902
			}]`),
			wantTemp: 0.0,
		},
		{
			testData: []byte(`[
				{
					"id":429009,
					"the_temp":34.5,
					"wind_speed":9.2608902
			}]`),
			wantTemp: 34.5,
		},
	}

	for _, iter := range temperatureRange {
		temperature, err := getCurrentTemperatureForCoordinates(iter.testData)
		if err == nil {
			if temperature != iter.wantTemp {
				t.Errorf("temperature mismatch, read %4.2f, want %4.2f", temperature, iter.wantTemp)
			}
		}
	}
}

func TestGetCityCoordinates(t *testing.T) {
	cityTestData := []byte(`[
    {
        "city": "New York", 
        "growth_from_2000_to_2013": "4.8%", 
        "latitude": 40.7127837, 
        "longitude": -74.0059413, 
        "population": "8405837", 
        "rank": "1", 
        "state": "New York"
    }, 
    {
        "city": "Los Angeles", 
        "growth_from_2000_to_2013": "4.8%", 
        "latitude": 34.0522342, 
        "longitude": -118.2436849, 
        "population": "3884307", 
        "rank": "2", 
        "state": "California"
	 },
	 {
        "city": "Chicago", 
        "growth_from_2000_to_2013": "-6.1%", 
        "latitude": 41.8781136, 
        "longitude": -87.6297982, 
        "population": "2718782", 
        "rank": "3", 
        "state": "Illinois"
    }
	 ]`)

	coordinatesRange := []struct {
		want coordinate
	}{
		{
			want: coordinate{ // New York
				Latitude:  40.7127837,
				Longitude: -74.0059413,
			},
		},
		{
			want: coordinate{ // Los Angeles
				Latitude:  34.0522342,
				Longitude: -118.2436849,
			},
		},
		{
			want: coordinate{ // Chicago
				Latitude:  41.8781136,
				Longitude: -87.6297982,
			},
		},
	}

	cityCoordinates := getCityCoordinates(cityTestData)

	for i, iter := range coordinatesRange {
		if iter.want != cityCoordinates[i] {
			t.Errorf("coordinates mismatch, read %.7f, want %.7f", cityCoordinates[i], iter.want)
		}
	}
}

func TestAverages(t *testing.T) {
	averagesRange := []struct {
		average []float64
		want    float64
	}{
		{
			average: []float64{
				27.35,
				33.405,
			},
			want: 30.3775,
		},
		{
			average: []float64{
				1.234,
				33.405,
				78.2,
				91.5,
			},
			want: 51.08475,
		},
		{
			average: []float64{
				3.0,
				3.0,
				3.0,
			},
			want: 3.0,
		},
	}

	for _, iter := range averagesRange {
		avg := getAverage(iter.average)
		if avg != iter.want {
			t.Errorf("average mismatch, read %.7f, want %.7f", avg, iter.want)
		}
	}
}

func TestGetCityWoeid(t *testing.T) {
	cityData := []byte(`[
				{
					"distance":1836,
					"title":"Santa Cruz",
					"location_type":"City",
					"woeid":2488853,
					"latt_long":"36.974018,-122.030952"
					}
				]`)
	want := int64(2488853)
	woeid := getCityWoeid(cityData)

	if woeid != want {
		t.Errorf("average mismatch, read %d, want %d", woeid, want)
	}
} 
