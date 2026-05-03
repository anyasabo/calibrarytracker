package scraper

// These types mirror the TypeScript definitions in src/lib/types/index.ts.
// If you change the JSON schema here, update the TypeScript types to match.

type DigitalAccess struct {
	HasOverdrive bool   `json:"hasOverdrive"`
	OverdriveURL string `json:"overdriveUrl"`
	OffersEcard  bool   `json:"offersEcard"`
	EcardURL     string `json:"ecardUrl"`
	EcardNotes   string `json:"ecardNotes"`
}

type LibrarySystem struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	County        string        `json:"county"`
	City          string        `json:"city"`
	Address       string        `json:"address"`
	Phone         string        `json:"phone"`
	Website       string        `json:"website"`
	Email         string        `json:"email"`
	CatalogURL    string        `json:"catalogUrl"`
	DirectorName  string        `json:"directorName"`
	CLSASystem    string        `json:"clsaSystem"`
	Lat           float64       `json:"lat"`
	Lng           float64       `json:"lng"`
	DigitalAccess DigitalAccess `json:"digitalAccess"`
	LastUpdated   string        `json:"lastUpdated"`
}

type Branch struct {
	ID          string  `json:"id"`
	SystemID    string  `json:"systemId"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zipCode"`
	County      string  `json:"county"`
	Phone       string  `json:"phone"`
	Website     string  `json:"website"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	OutletType  string  `json:"outletType"`
	Status      string  `json:"status"`
	LastUpdated string  `json:"lastUpdated"`
}
