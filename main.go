package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type (
	Weather struct {
		Location struct {
			Name    string `json:"name"`
			Country string `json:"country"`
		} `json:"location"`
		Current struct {
			LastUpdatedEpoch int64   `json:"last_updated_epoch"`
			Temperature      float64 `json:"temp_c"`
			Condition        struct {
				Text string `json:"text"`
			} `json:"condition"`
			WindSpeed    float64 `json:"wind_kph"`
			Gusts        float64 `json:"gust_kph"`
			ChangeofRain int64   `json:"chance_of_rain"`
		} `json:"current"`
		Forecast struct {
			ForecastDay []struct {
				Hours []struct {
					DateEpoch   int64   `json:"time_epoch"`
					Temperature float64 `json:"temp_c"`
					Condition   struct {
						Text string `json:"text"`
					} `json:"condition"`
					WindSpeed    float64 `json:"wind_kph"`
					Gusts        float64 `json:"gust_kph"`
					ChangeofRain int64   `json:"chance_of_rain"`
				} `json:"hour"`
			} `json:"forecastday"`
		} `json:"forecast"`
	}
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Read environment variables
	accesstoken := os.Getenv("WEATHER_ACCESS_TOKEN")
	q := os.Getenv("DEFAULT_LOCATION")

	// Deterime if there is a location
	// specified in the command line
	if len(os.Args) > 1 {
		q = os.Args[1]
	}

	// Construct the URL
	url := "https://api.weatherapi.com/v1/forecast.json?q=" + q + "&days=1&key=" + accesstoken

	// Call the API
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// If Status is not OK 200, panic
	if res.StatusCode != 200 {
		panic("Weather Api not available")
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the json into the provide
	// struct reference
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	// Extract data from the struct
	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hours

	// Build header message
	message := fmt.Sprintf("\n\n%s, %s\nCurrent Conditions\n",
		location.Name, location.Country)

	fmt.Println(message)

	// Build and display current conditions
	message = fmt.Sprintf("%.0f°c, %s. %d%% chance of rain\nwind %.0fkm/h, gusts %.0fkm/h\n\n",
		current.Temperature, current.Condition.Text, current.ChangeofRain, current.WindSpeed, current.Gusts)
	color.Cyan(message)

	fmt.Println("Forecast:")

	// Get the hourly forcasts and
	// construct an output message
	for _, hour := range hours {
		date := time.Unix(hour.DateEpoch, 0)

		// If the hourly forcast is in the past
		// ignore it and continue along
		if date.Before(time.Now()) {
			continue
		}

		message = fmt.Sprintf("%s | %.0f°c | %s | %d%% rain | wind %.0f(%.0f) km/h",
			date.Format("15:04"), hour.Temperature, hour.Condition.Text, hour.ChangeofRain, hour.WindSpeed, hour.Gusts)

		// Chance of rain and Wind gusts
		// will change the color.
		if hour.ChangeofRain < 50 {
			if hour.Gusts < 45 {
				color.Green(message)
			} else {
				color.Yellow(message)
			}
		} else {
			if hour.Gusts < 45 {
				color.Cyan(message)
			} else {
				color.Red(message)
			}
		}
	}

}
