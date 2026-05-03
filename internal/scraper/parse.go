package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	adminHQURL = "https://ca.countingopinions.com/pireports/report.php?b93196abc5834dc7148fb08c850be3f0&live"
	branchURL  = "https://ca.countingopinions.com/pireports/report.php?c5e42ef64d3cd9b3c91539fa4a814312&live"
)

func fetchHTML(url string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetching %s: status %d", url, resp.StatusCode)
	}
	return resp.Body, nil
}

// ParseAdminHQ fetches and parses the administrative headquarters directory.
func ParseAdminHQ() ([]LibrarySystem, error) {
	body, err := fetchHTML(adminHQURL)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseAdminHQHTML(body)
}

func parseAdminHQHTML(r io.Reader) ([]LibrarySystem, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("parsing admin HQ HTML: %w", err)
	}

	now := time.Now().UTC().Format(time.DateOnly)
	var systems []LibrarySystem
	seen := make(map[string]bool)

	// Detect column layout from the header row. We find the first <tr> and
	// iterate its cells. This works with or without <thead>.
	colMap := make(map[string]int)
	headerRow := doc.Find("table tr").First()
	headerRow.Find("td, th").Each(func(i int, s *goquery.Selection) {
		text := strings.ToLower(cleanText(s.Text()))
		switch {
		case strings.Contains(text, "location"):
			colMap["name"] = i
		case strings.Contains(text, "street address") && !strings.Contains(text, "mailing"):
			if _, exists := colMap["address"]; !exists {
				colMap["address"] = i
			}
		case text == "1.11 city" || (strings.Contains(text, "city") && !strings.Contains(text, "mailing")):
			if _, exists := colMap["city"]; !exists {
				colMap["city"] = i
			}
		case strings.Contains(text, "zip") && !strings.Contains(text, "+4") && !strings.Contains(text, "mailing"):
			if _, exists := colMap["zip"]; !exists {
				colMap["zip"] = i
			}
		case strings.Contains(text, "phone") && strings.Contains(text, "admin"):
			colMap["phone"] = i
		case strings.Contains(text, "email") || strings.Contains(text, "contact"):
			colMap["email"] = i
		case strings.Contains(text, "web address") || strings.Contains(text, "web site"):
			colMap["website"] = i
		case strings.Contains(text, "first name"):
			colMap["dirFirst"] = i
		case strings.Contains(text, "last name"):
			colMap["dirLast"] = i
		}
	})

	// Fallback column indices if header detection fails.
	// Based on observed HTML structure from countingopinions.com:
	// 0: Name, 1: Address, 2: City, 3: Zip, 4: Zip+4,
	// 5: Admin Phone, 6: Ref Phone, 7: TDD,
	// 8: Email, 9: Website, 10: Dir First, 11: Dir Last, 12: Dir Title
	defaults := map[string]int{
		"name": 0, "address": 1, "city": 2, "zip": 3,
		"phone": 5, "email": 8, "website": 9,
		"dirFirst": 10, "dirLast": 11,
	}
	for k, v := range defaults {
		if _, exists := colMap[k]; !exists {
			colMap[k] = v
		}
	}

	// Process all table rows. We use "table tr" to cover tables with or without
	// <thead>/<tbody>. The deduplication via `seen` map and the header-row skip
	// below handle any double-matching.
	doc.Find("table tr").Each(func(idx int, row *goquery.Selection) {
		cells := row.Find("td, th")
		if cells.Length() < 10 {
			return
		}

		col := func(key string) string {
			i, ok := colMap[key]
			if !ok || i >= cells.Length() {
				return ""
			}
			return cleanText(cells.Eq(i).Text())
		}

		name := col("name")
		if name == "" {
			return
		}
		// Skip header rows
		lower := strings.ToLower(name)
		if strings.Contains(lower, "location") || strings.Contains(lower, "10.5") {
			return
		}

		id := Slugify(name)
		if seen[id] {
			return
		}
		seen[id] = true

		email := col("email")
		website := col("website")
		if strings.HasPrefix(email, "http") && !strings.HasPrefix(website, "http") {
			email, website = website, email
		}

		dirName := ""
		firstName := col("dirFirst")
		lastName := col("dirLast")
		if firstName != "" || lastName != "" {
			dirName = strings.TrimSpace(firstName + " " + lastName)
		}

		sys := LibrarySystem{
			ID:            id,
			Name:          name,
			Address:       col("address"),
			City:          col("city"),
			Phone:         col("phone"),
			Website:       website,
			Email:         email,
			DirectorName:  dirName,
			LastUpdated:   now,
			DigitalAccess: DigitalAccess{},
		}

		systems = append(systems, sys)
	})

	return systems, nil
}

// ParseBranches fetches and parses the branch/outlet directory.
func ParseBranches() ([]Branch, error) {
	body, err := fetchHTML(branchURL)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseBranchesHTML(body)
}

func parseBranchesHTML(r io.Reader) ([]Branch, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("parsing branches HTML: %w", err)
	}

	// Branch report column layout (from countingopinions.com):
	// 0: Location (short name)
	// 1: 10.5 Name (full name)
	// 2: 10.6 Street Address
	// 3: 10.7 City
	// 4: 10.8 Zip Code
	// 5: 10.9 Zip+4
	// 6: 10.10 Mailing Street
	// 7: 10.11 Mailing City
	// 8: 10.12 Mailing Zip
	// 9: 10.13 Mailing Zip+4
	// 10: 10.14 County
	// 11: 10.15 Phone
	// 12: Latitude
	// 13: Longitude
	// 14: 10.16 Outlet Type
	// 15: Main Library Name
	// 16: Status
	// 17: Law Library

	now := time.Now().UTC().Format(time.DateOnly)
	var branches []Branch
	seen := make(map[string]bool)

	doc.Find("table tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td, th")
		if cells.Length() < 16 {
			return
		}

		// Prefer the full name (col 1), fall back to location (col 0)
		name := cleanText(cells.Eq(1).Text())
		if name == "" {
			name = cleanText(cells.Eq(0).Text())
		}
		if name == "" {
			return
		}
		// Skip header rows
		lowerName := strings.ToLower(name)
		if strings.Contains(lowerName, "10.5 name") || strings.Contains(lowerName, "location") {
			return
		}

		mainLibName := cleanText(cells.Eq(15).Text())
		systemID := Slugify(mainLibName)

		id := Slugify(name)
		key := systemID + "/" + id
		if seen[key] {
			return
		}
		seen[key] = true

		lat, _ := strconv.ParseFloat(cleanText(cells.Eq(12).Text()), 64)
		lng, _ := strconv.ParseFloat(cleanText(cells.Eq(13).Text()), 64)

		outletType := parseOutletType(cleanText(cells.Eq(14).Text()))
		county := cleanCounty(cleanText(cells.Eq(10).Text()))

		b := Branch{
			ID:          id,
			SystemID:    systemID,
			Name:        name,
			Address:     cleanText(cells.Eq(2).Text()),
			City:        cleanText(cells.Eq(3).Text()),
			ZipCode:     cleanText(cells.Eq(4).Text()),
			County:      county,
			Phone:       cleanText(cells.Eq(11).Text()),
			Lat:         lat,
			Lng:         lng,
			OutletType:  outletType,
			Status:      "open",
			LastUpdated: now,
		}

		branches = append(branches, b)
	})

	return branches, nil
}

func parseOutletType(raw string) string {
	lower := strings.ToLower(raw)
	if strings.Contains(lower, "central") {
		return "central"
	}
	if strings.Contains(lower, "bookmobile") {
		return "bookmobile"
	}
	if strings.Contains(lower, "branch") {
		return "branch"
	}
	return "other"
}

// cleanCounty removes the leading numeric code from county strings
// like "36San Diego" -> "San Diego"
func cleanCounty(raw string) string {
	for i, c := range raw {
		if c < '0' || c > '9' {
			return raw[i:]
		}
	}
	return raw
}

func cleanText(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	// Collapse multiple spaces
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	if s == "-1" || s == "N/A" || s == "n/a" || s == "(___) ___-____" {
		return ""
	}
	return s
}
