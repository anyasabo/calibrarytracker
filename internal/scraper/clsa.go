package scraper

import (
	"encoding/json"
	"fmt"
	"os"
)

type cooperativeJSON struct {
	ID        string `json:"id"`
	Subgroups []struct {
		ID              string   `json:"id"`
		MemberSystemIDs []string `json:"memberSystemIds"`
	} `json:"subgroups"`
}

// PopulateCLSA reads cooperatives.json and sets CLSASystem on each
// matching LibrarySystem to the "{cooperative}-{subgroup}" slug.
func PopulateCLSA(path string, systems []LibrarySystem) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	var coops []cooperativeJSON
	if err := json.NewDecoder(f).Decode(&coops); err != nil {
		return fmt.Errorf("decoding %s: %w", path, err)
	}

	// Build system ID -> CLSA subgroup ID mapping
	membership := make(map[string]string)
	for _, coop := range coops {
		for _, sg := range coop.Subgroups {
			for _, sysID := range sg.MemberSystemIDs {
				membership[sysID] = sg.ID
			}
		}
	}

	for i := range systems {
		if sgID, ok := membership[systems[i].ID]; ok {
			systems[i].CLSASystem = sgID
		}
	}

	return nil
}
