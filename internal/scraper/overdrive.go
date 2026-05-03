package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const overdriveAPIBase = "https://thunder.api.overdrive.com/v2/libraries"

type odLibrary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Links struct {
		LibraryHome struct {
			Href string `json:"href"`
		} `json:"libraryHome"`
	} `json:"links"`
}

type odResponse struct {
	Items      []odLibrary `json:"items"`
	TotalItems int         `json:"totalItems"`
	Links      struct {
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
	} `json:"links"`
}

// FetchOverdriveLibraries paginates the Thunder API and returns all libraries.
func FetchOverdriveLibraries() ([]odLibrary, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	var all []odLibrary
	nextURL := overdriveAPIBase + "?limit=100"

	for nextURL != "" {
		req, err := http.NewRequest("GET", nextURL, nil)
		if err != nil {
			return nil, fmt.Errorf("building request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching %s: %w", nextURL, err)
		}

		var page odResponse
		decErr := json.NewDecoder(resp.Body).Decode(&page)
		resp.Body.Close()
		if decErr != nil {
			return nil, fmt.Errorf("decoding response: %w", decErr)
		}

		all = append(all, page.Items...)
		nextURL = page.Links.Next.Href
	}

	return all, nil
}

func normalizeDomain(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "http://" + rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := strings.ToLower(parsed.Hostname())
	host = strings.TrimPrefix(host, "www.")
	return host
}

// MatchOverdriveToSystems matches OverDrive libraries to our systems by
// website domain and name, and writes the results into digital-access.json.
// Preserves existing hand-curated eCard fields.
func MatchOverdriveToSystems(odLibs []odLibrary, systems []LibrarySystem, digitalAccessPath string) error {
	// Load existing digital-access.json to preserve eCard data
	existing := make(map[string]digitalAccessEntry)
	if f, err := os.Open(digitalAccessPath); err == nil {
		json.NewDecoder(f).Decode(&existing)
		f.Close()
	}

	// Build lookup maps from our systems
	domainToID := make(map[string]string)
	nameToID := make(map[string]string)
	for _, s := range systems {
		d := normalizeDomain(s.Website)
		if d != "" {
			domainToID[d] = s.ID
		}
		nameToID[strings.ToLower(s.Name)] = s.ID
	}

	matched := 0
	for _, lib := range odLibs {
		homeURL := lib.Links.LibraryHome.Href
		odDomain := normalizeDomain(homeURL)

		sysID := ""
		if odDomain != "" {
			sysID = domainToID[odDomain]
		}
		if sysID == "" {
			sysID = nameToID[strings.ToLower(lib.Name)]
		}
		if sysID == "" {
			continue
		}

		entry := existing[sysID]
		entry.HasOverdrive = true
		odURL := fmt.Sprintf("https://%s.overdrive.com", lib.ID)
		entry.OverdriveURL = odURL
		existing[sysID] = entry
		matched++
	}

	fmt.Printf("  Matched %d OverDrive libraries to systems\n", matched)

	// Write back
	f, err := os.Create(digitalAccessPath)
	if err != nil {
		return fmt.Errorf("creating %s: %w", digitalAccessPath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(existing)
}
