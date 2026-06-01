package sources

import (
	"context"
	"fmt"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
)

// Weather is current conditions for a location.
type Weather struct {
	TemperatureC float64   `json:"temperature_c"`
	Humidity     float64   `json:"humidity_pct"`
	WindSpeedKmh float64   `json:"wind_speed_kmh"`
	WeatherCode  int       `json:"weather_code"`
	ObservedAt   time.Time `json:"observed_at"`
}

// FetchWeather gets current conditions from Open-Meteo.
func FetchWeather(ctx context.Context, lat, lon float64) (Weather, domain.Source, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
		"&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code", lat, lon)
	var raw struct {
		Current struct {
			Time        string  `json:"time"`
			Temperature float64 `json:"temperature_2m"`
			Humidity    float64 `json:"relative_humidity_2m"`
			WindSpeed   float64 `json:"wind_speed_10m"`
			WeatherCode int     `json:"weather_code"`
		} `json:"current"`
	}
	if err := getJSON(ctx, url, &raw); err != nil {
		return Weather{}, domain.Source{}, err
	}
	obs, _ := time.Parse("2006-01-02T15:04", raw.Current.Time)
	w := Weather{
		TemperatureC: raw.Current.Temperature,
		Humidity:     raw.Current.Humidity,
		WindSpeedKmh: raw.Current.WindSpeed,
		WeatherCode:  raw.Current.WeatherCode,
		ObservedAt:   obs.UTC(),
	}
	return w, openMeteoSource(), nil
}

// AirQuality is current air-quality for a location.
type AirQuality struct {
	USAQI      int       `json:"us_aqi"`
	PM10       float64   `json:"pm10"`
	PM25       float64   `json:"pm2_5"`
	ObservedAt time.Time `json:"observed_at"`
}

// FetchAirQuality gets current AQI from the Open-Meteo air-quality API.
func FetchAirQuality(ctx context.Context, lat, lon float64) (AirQuality, domain.Source, error) {
	url := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%.4f&longitude=%.4f"+
		"&current=us_aqi,pm10,pm2_5", lat, lon)
	var raw struct {
		Current struct {
			Time  string  `json:"time"`
			USAQI float64 `json:"us_aqi"`
			PM10  float64 `json:"pm10"`
			PM25  float64 `json:"pm2_5"`
		} `json:"current"`
	}
	if err := getJSON(ctx, url, &raw); err != nil {
		return AirQuality{}, domain.Source{}, err
	}
	obs, _ := time.Parse("2006-01-02T15:04", raw.Current.Time)
	aq := AirQuality{
		USAQI:      int(raw.Current.USAQI),
		PM10:       raw.Current.PM10,
		PM25:       raw.Current.PM25,
		ObservedAt: obs.UTC(),
	}
	return aq, openMeteoSource(), nil
}

func openMeteoSource() domain.Source {
	return domain.Source{Name: "Open-Meteo", URL: "https://open-meteo.com", FetchedAt: time.Now().UTC()}
}
