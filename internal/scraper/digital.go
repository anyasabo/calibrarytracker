package scraper

import (
	"encoding/json"
	"fmt"
	"os"
)

type digitalAccessEntry struct {
	HasOverdrive bool   `json:"hasOverdrive"`
	OverdriveURL string `json:"overdriveUrl"`
	OffersEcard  bool   `json:"offersEcard"`
	EcardURL     string `json:"ecardUrl"`
	EcardNotes   string `json:"ecardNotes"`
}

// LoadDigitalAccess reads digital-access.json and merges it into systems.
func LoadDigitalAccess(path string, systems []LibrarySystem) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	var entries map[string]digitalAccessEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return fmt.Errorf("decoding %s: %w", path, err)
	}

	for i := range systems {
		entry, ok := entries[systems[i].ID]
		if !ok {
			continue
		}
		systems[i].DigitalAccess = DigitalAccess{
			HasOverdrive: entry.HasOverdrive,
			OverdriveURL: entry.OverdriveURL,
			OffersEcard:  entry.OffersEcard,
			EcardURL:     entry.EcardURL,
			EcardNotes:   entry.EcardNotes,
		}
	}

	return nil
}
