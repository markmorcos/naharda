package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/markmorcos/naharda/api/internal/calendar"
	"github.com/markmorcos/naharda/api/internal/domain"
	"github.com/markmorcos/naharda/api/internal/http/respond"
	"github.com/markmorcos/naharda/api/internal/sources"
)

// Cities returns the 13 canonical cities (§4). Long-cached.
func (h *Handlers) Cities(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, r, 86400, domain.Cities, domain.Meta{
		Sources:     []domain.Source{{Name: "Naharda", FetchedAt: time.Now().UTC()}},
		Attribution: "Cities curated by Naharda. Free tier non-commercial.",
	})
}

// Calendar converts Gregorian ↔ Hijri locally (§4). Optional ?date=YYYY-MM-DD.
func (h *Handlers) Calendar(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UTC()
	if d := r.URL.Query().Get("date"); d != "" {
		parsed, err := time.Parse("2006-01-02", d)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid_date", "Use date=YYYY-MM-DD.", 0)
			return
		}
		t = parsed
	}
	conv, err := calendar.Convert(t)
	if err != nil {
		respond.Error(w, http.StatusUnprocessableEntity, "out_of_range", err.Error(), 0)
		return
	}
	respond.JSON(w, r, 86400, conv, domain.Meta{
		Sources:     []domain.Source{{Name: "Naharda local compute (Umm al-Qura)", FetchedAt: time.Now().UTC()}},
		Attribution: "Hijri date computed locally by Naharda. Free tier non-commercial.",
	})
}

// PrayerTimes returns salah times for a city (§4). Today 1h; past immutable.
func (h *Handlers) PrayerTimes(w http.ResponseWriter, r *http.Request) {
	city, ok := domain.CityBySlug(chi.URLParam(r, "city"))
	if !ok {
		respond.Error(w, http.StatusNotFound, "unknown_city", "Unknown city.", 0)
		return
	}
	date := time.Now().UTC()
	past := false
	if d := r.URL.Query().Get("date"); d != "" {
		parsed, err := time.Parse("2006-01-02", d)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid_date", "Use date=YYYY-MM-DD.", 0)
			return
		}
		date = parsed
		past = parsed.Before(time.Now().UTC().Truncate(24 * time.Hour))
	}
	pt, src, err := sources.FetchPrayerTimes(r.Context(), city.Lat, city.Lon, date, sources.MethodEgypt)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, "upstream_error", "Prayer-times source unavailable.", 30)
		return
	}
	maxAge := 3600
	if past {
		maxAge = 31536000 // immutable past (§9.2)
	}
	respond.JSON(w, r, maxAge, map[string]any{"city": city.Slug, "prayer_times": pt}, domain.Meta{
		Sources:     []domain.Source{src},
		Attribution: "Prayer times via Aladhan. Free tier non-commercial.",
	})
}

// Weather returns current conditions for a city (§4). 10-minute cache.
func (h *Handlers) Weather(w http.ResponseWriter, r *http.Request) {
	city, ok := domain.CityBySlug(chi.URLParam(r, "city"))
	if !ok {
		respond.Error(w, http.StatusNotFound, "unknown_city", "Unknown city.", 0)
		return
	}
	wx, src, err := sources.FetchWeather(r.Context(), city.Lat, city.Lon)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, "upstream_error", "Weather source unavailable.", 30)
		return
	}
	respond.JSON(w, r, 600, map[string]any{"city": city.Slug, "weather": wx}, domain.Meta{
		Sources:     []domain.Source{src},
		Attribution: "Weather via Open-Meteo. Free tier non-commercial.",
	})
}

// AirQuality returns current AQI for a city (§4). 10-minute cache.
func (h *Handlers) AirQuality(w http.ResponseWriter, r *http.Request) {
	city, ok := domain.CityBySlug(chi.URLParam(r, "city"))
	if !ok {
		respond.Error(w, http.StatusNotFound, "unknown_city", "Unknown city.", 0)
		return
	}
	aq, src, err := sources.FetchAirQuality(r.Context(), city.Lat, city.Lon)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, "upstream_error", "Air-quality source unavailable.", 30)
		return
	}
	respond.JSON(w, r, 600, map[string]any{"city": city.Slug, "aqi": aq}, domain.Meta{
		Sources:     []domain.Source{src},
		Attribution: "Air quality via Open-Meteo. Free tier non-commercial.",
	})
}

// Fuel returns official pump prices (§4), preferring manual overrides. Day cache.
func (h *Handlers) Fuel(w http.ResponseWriter, r *http.Request) {
	prices := make([]domain.FuelPrice, len(domain.DefaultFuelPrices))
	copy(prices, domain.DefaultFuelPrices)
	for i, p := range prices {
		if v, ok, err := h.store.ActiveOverride(r.Context(), "fuel", p.Product); err == nil && ok {
			prices[i].PriceEGP = v
		}
	}
	respond.JSON(w, r, 86400, prices, domain.Meta{
		Sources:     []domain.Source{{Name: "Egyptian Ministry of Petroleum / EGPC", FetchedAt: time.Now().UTC()}},
		Attribution: "Fuel prices: Egyptian Ministry of Petroleum / EGPC. Free tier non-commercial.",
	})
}
