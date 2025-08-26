package sun

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	cityFlag string
	daysFlag int
)

// Run is the entry point of the application. It loads environment variables,
// parses command-line flags, constructs a URL to call a weather API, and processes
// the response to display current weather conditions and a forecast.
//
// It handles errors by reporting them to stderr and returning early (no panics).
func Run() {
	// Load environment variables from .env file (optional)
	setupEnvironment()

	// Get the number of days to forecast
	noOfDaysStr := os.Getenv("NO_OF_DAYS")
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

	// Call the API with context and timeout
	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "build request error: %v\n", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		return
	}
	defer func() {
		if cErr := res.Body.Close(); cErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error closing response body: %v\n", cErr)
		}
	}()

	// If Status is not OK, report and return
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		_, _ = fmt.Fprintf(os.Stderr, "weather api: status %d: %s\n", res.StatusCode, string(b))
		return
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "read body error: %v\n", err)
		return
	}

	// Unmarshal the json into the provided struct reference
	var weather Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "json decode error: %v\n", err)
		return
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
		date := time.Unix(fday.DateEpoch, 0).Format("2006-01-02")

		// get the name of the day for date
		dayName := time.Unix(fday.DateEpoch, 0).Weekday().String()

		fmt.Printf("%s (%s)\n", date, dayName)
		// Get the hourly forecasts and construct an output message
		for _, hour := range hours {
			hourDate := time.Unix(hour.DateEpoch, 0)
			if hourDate.Before(time.Now()) {
				continue
			}

			message = fmt.Sprintf("%s | %2.0f°c | %-20s | %4d%% rain | wind %3.0f(%2.0f) km/h",
				hourDate.Format("15:04"), hour.Temperature, hour.Condition.Text, hour.ChanceOfRain, hour.WindSpeed, hour.Gusts)

			// Chance of rain and Wind gusts will change the color.
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
	// .env is optional in many environments; log a warning instead of panicking
	if err := godotenv.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: .env not loaded: %v\n", err)
	}
}
