// Package calendar converts Gregorian ↔ Hijri locally with no external
// dependency (project.md §4). Uses the Umm al-Qura civil Islamic calendar.
package calendar

import (
	"fmt"
	"time"

	hijri "github.com/hablullah/go-hijri"
)

var hijriMonths = []string{
	"", "Muharram", "Safar", "Rabi al-Awwal", "Rabi al-Thani",
	"Jumada al-Awwal", "Jumada al-Thani", "Rajab", "Sha'ban",
	"Ramadan", "Shawwal", "Dhu al-Qa'dah", "Dhu al-Hijjah",
}

// Conversion is a Gregorian/Hijri date pair.
type Conversion struct {
	Gregorian      string `json:"gregorian"` // YYYY-MM-DD
	Hijri          string `json:"hijri"`     // YYYY-MM-DD
	HijriDay       int    `json:"hijri_day"`
	HijriMonth     int    `json:"hijri_month"`
	HijriMonthName string `json:"hijri_month_name"`
	HijriYear      int    `json:"hijri_year"`
	Weekday        string `json:"weekday"`
}

// Convert returns the Hijri date for a Gregorian instant.
func Convert(t time.Time) (Conversion, error) {
	uq, err := hijri.CreateUmmAlQuraDate(t)
	if err != nil {
		return Conversion{}, err
	}
	month := ""
	if uq.Month >= 1 && int(uq.Month) < len(hijriMonths) {
		month = hijriMonths[uq.Month]
	}
	return Conversion{
		Gregorian:      t.Format("2006-01-02"),
		Hijri:          fmt.Sprintf("%04d-%02d-%02d", uq.Year, uq.Month, uq.Day),
		HijriDay:       int(uq.Day),
		HijriMonth:     int(uq.Month),
		HijriMonthName: month,
		HijriYear:      int(uq.Year),
		Weekday:        t.Weekday().String(),
	}, nil
}
