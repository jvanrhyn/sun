
# Weather Forecast Application 

A simple Go application that fetches and displays the current weather conditions and forecast for a specified city using the WeatherAPI.

## Features

- Fetches current weather conditions.
- Displays a weather forecast for a specified number of days.
- Color-coded output based on weather conditions.

## Getting Started

### Prerequisites

- Go 1.16 or later
- A WeatherAPI access token
- A `.env` file with the following variables:
- Obtain a free API key from [WeatherAPI](https://www.weatherapi.com/). 

```
WEATHER_ACCESS_TOKEN=your_weather_api_token
DEFAULT_LOCATION=your_default_city
NO_OF_DAYS=number_of_days_to_forecast
```

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/jvanrhyn/sun.git
   cd sun
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

### Usage

1. Run the application:
   ```sh
   go run main.go
   ```

2. You can specify the city and number of days to forecast using command-line flags:
   ```sh
   go run main.go -city="New York" -days=3
   ```

### Example Output

```
New York, United States
Current Conditions
25°c, Partly cloudy. 10% chance of rain
wind 15km/h, gusts 20km/h

Forecast:
2023-10-01
15:00 | 22°c | Clear | 0% rain | wind 10(15) km/h
16:00 | 23°c | Clear | 0% rain | wind 12(18) km/h
...
```

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
