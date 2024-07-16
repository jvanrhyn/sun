package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var (
	cityFlag string
	daysFlag int
)

// main is the entry point of the application. It loads environment variables,
// parses command-line flags, constructs a URL to call a weather API, and processes
// the response to display current weather conditions and a forecast.
//
// The function does not take any parameters and does not return any values.
func main() {

	// Load environment variables from .env file
	setupEnvironment()

	// Get the number of days to forecast
	var noOfDaysStr = os.Getenv("NO_OF_DAYS")
	noOfDays, err := strconv.Atoi(noOfDaysStr)
	if err != nil {
		noOfDays = 1
	}

	// Parse the command-line flags
	flag.IntVar(&daysFlag, "days", noOfDays, "Number of days to forecast")
	flag.StringVar(&cityFlag, "city", os.Getenv("DEFAULT_LOCATION"), "Enter the name of the city")
	flag.Parse()

	// Retrieve the access token from the environment
	token := os.Getenv("WEATHER_ACCESS_TOKEN")

	// Construct the URL
	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?q=%s&days=%d&key=%s",
		cityFlag, daysFlag, token)

	// Call the API
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error closing response body: %v\n", cerr)
		}
	}()

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
	location, current, forecastDay := weather.Location, weather.Current, weather.Forecast.ForecastDay

	// Build header message
	message := fmt.Sprintf("\n\n%s, %s\nCurrent Conditions\n",
		location.Name, location.Country)
	fmt.Println(message)

	// Build and display current conditions
	message = fmt.Sprintf("%.0f°c, %s. %d%% chance of rain\nwind %.0fkm/h, gusts %.0fkm/h\n\n",
		current.Temperature, current.Condition.Text, current.ChanceOfRain, current.WindSpeed, current.Gusts)
	color.Cyan(message)

	fmt.Println("Forecast:")

	for _, fday := range forecastDay {
		hours := fday.Hours

		fdate := time.Unix(fday.DateEpoch, 0).Format("2006-01-02")

		fmt.Println(fdate)
		// Get the hourly forecasts and
		// construct an output message
		for _, hour := range hours {
			date := time.Unix(hour.DateEpoch, 0)

			// If the hourly forecast is in the past
			// ignore it and continue along
			if date.Before(time.Now()) {
				continue
			}

			message = fmt.Sprintf("%s | %.0f°c | %s | %d%% rain | wind %.0f(%.0f) km/h",
				date.Format("15:04"), hour.Temperature, hour.Condition.Text, hour.ChanceOfRain, hour.WindSpeed, hour.Gusts)

			// Chance of rain and Wind gusts
			// will change the color.
			if hour.ChanceOfRain < 50 {
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
}

func setupEnvironment() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
