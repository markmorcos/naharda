package domain

// City is a canonical Naharda city with hardcoded coordinates (project.md §4).
type City struct {
	Slug string  `json:"slug"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// Cities is the v1 canonical, hardcoded list (project.md §4).
var Cities = []City{
	{"cairo", "Cairo", 30.0444, 31.2357},
	{"giza", "Giza", 30.0131, 31.2089},
	{"alexandria", "Alexandria", 31.2001, 29.9187},
	{"hurghada", "Hurghada", 27.2579, 33.8116},
	{"sharm-el-sheikh", "Sharm El Sheikh", 27.9158, 34.3300},
	{"aswan", "Aswan", 24.0889, 32.8998},
	{"luxor", "Luxor", 25.6872, 32.6396},
	{"mansoura", "Mansoura", 31.0409, 31.3785},
	{"tanta", "Tanta", 30.7865, 31.0004},
	{"asyut", "Asyut", 27.1809, 31.1837},
	{"port-said", "Port Said", 31.2653, 32.3019},
	{"suez", "Suez", 29.9668, 32.5498},
	{"ismailia", "Ismailia", 30.6043, 32.2723},
}

// CityBySlug returns the city with the given slug, if present.
func CityBySlug(slug string) (City, bool) {
	for _, c := range Cities {
		if c.Slug == slug {
			return c, true
		}
	}
	return City{}, false
}
