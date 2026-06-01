package sources

import (
	"context"
	"fmt"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
)

// MethodEgypt is the Egyptian General Authority of Survey method code (Aladhan).
const MethodEgypt = 5

// PrayerTimes holds the day's salah times with the calculation method (§4).
type PrayerTimes struct {
	Date    string            `json:"date"`   // YYYY-MM-DD
	Method  string            `json:"method"` // explicit method name
	Timings map[string]string `json:"timings"`
}

// FetchPrayerTimes gets prayer times for a location and date from Aladhan.
func FetchPrayerTimes(ctx context.Context, lat, lon float64, date time.Time, method int) (PrayerTimes, domain.Source, error) {
	url := fmt.Sprintf("https://api.aladhan.com/v1/timings/%s?latitude=%.4f&longitude=%.4f&method=%d",
		date.Format("02-01-2006"), lat, lon, method)
	var raw struct {
		Data struct {
			Timings map[string]string `json:"timings"`
			Meta    struct {
				Method struct {
					Name string `json:"name"`
				} `json:"method"`
			} `json:"meta"`
		} `json:"data"`
	}
	if err := getJSON(ctx, url, &raw); err != nil {
		return PrayerTimes{}, domain.Source{}, err
	}
	keep := []string{"Fajr", "Sunrise", "Dhuhr", "Asr", "Maghrib", "Isha"}
	timings := make(map[string]string, len(keep))
	for _, k := range keep {
		if v, ok := raw.Data.Timings[k]; ok {
			timings[k] = v
		}
	}
	pt := PrayerTimes{
		Date:    date.Format("2006-01-02"),
		Method:  raw.Data.Meta.Method.Name,
		Timings: timings,
	}
	src := domain.Source{Name: "Aladhan", URL: "https://aladhan.com", FetchedAt: time.Now().UTC()}
	return pt, src, nil
}
