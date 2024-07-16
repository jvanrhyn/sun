package main

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
			ChanceOfRain int64   `json:"chance_of_rain"`
		} `json:"current"`
		Forecast struct {
			ForecastDay []struct {
				DateEpoch int64 `json:"date_epoch"`
				Hours     []struct {
					DateEpoch   int64   `json:"time_epoch"`
					Temperature float64 `json:"temp_c"`
					Condition   struct {
						Text string `json:"text"`
					} `json:"condition"`
					WindSpeed    float64 `json:"wind_kph"`
					Gusts        float64 `json:"gust_kph"`
					ChanceOfRain int64   `json:"chance_of_rain"`
				} `json:"hour"`
			} `json:"forecastday"`
		} `json:"forecast"`
	}
)
