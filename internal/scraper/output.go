package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// WriteJSON writes data to a JSON file with consistent formatting.
// The output is sorted by ID for stable diffs across scraper runs.
func WriteJSON(outputDir string, filename string, data interface{}) error {
	path := filepath.Join(outputDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// SortSystems sorts systems by ID for deterministic output.
func SortSystems(systems []LibrarySystem) {
	sort.Slice(systems, func(i, j int) bool {
		return systems[i].ID < systems[j].ID
	})
}

// SortBranches sorts branches by system ID then branch ID.
func SortBranches(branches []Branch) {
	sort.Slice(branches, func(i, j int) bool {
		if branches[i].SystemID != branches[j].SystemID {
			return branches[i].SystemID < branches[j].SystemID
		}
		return branches[i].ID < branches[j].ID
	})
}

// Known mismatches between branch report "Main Library Name" and admin HQ
// "Location" name. The branch report uses an older or abbreviated name;
// we map it to the admin HQ slug.
var branchSystemAliases = map[string]string{
	"glendale-public-library": "glendale-library-arts-culture",
	"goleta-valley-library":   "goleta-santa-ynez-valley-libraries",
	"santa-ana":               "santa-ana-public-library",
}

// CrossReference enriches systems with county data from branches
// and logs any branches that reference unknown systems.
func CrossReference(systems []LibrarySystem, branches []Branch) {
	systemMap := make(map[string]*LibrarySystem, len(systems))
	for i := range systems {
		systemMap[systems[i].ID] = &systems[i]
	}

	for i := range branches {
		// Resolve known aliases
		if alias, ok := branchSystemAliases[branches[i].SystemID]; ok {
			branches[i].SystemID = alias
		}

		sys, ok := systemMap[branches[i].SystemID]
		if !ok {
			if branches[i].SystemID != "" {
				fmt.Fprintf(os.Stderr, "warning: branch %q references unknown system %q\n", branches[i].Name, branches[i].SystemID)
			}
			continue
		}
		if sys.County == "" && branches[i].County != "" {
			sys.County = branches[i].County
		}
		// Prefer central branch coordinates; fall back to any branch
		if branches[i].Lat != 0 && branches[i].Lng != 0 {
			if sys.Lat == 0 && sys.Lng == 0 {
				sys.Lat = branches[i].Lat
				sys.Lng = branches[i].Lng
			} else if branches[i].OutletType == "central" {
				sys.Lat = branches[i].Lat
				sys.Lng = branches[i].Lng
			}
		}
	}
}
