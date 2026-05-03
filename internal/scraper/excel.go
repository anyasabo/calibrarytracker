package scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// SpreadsheetML types for parsing the Excel XML export from countingopinions.com
type xmlWorkbook struct {
	XMLName xml.Name       `xml:"Workbook"`
	Sheets  []xmlWorksheet `xml:"Worksheet"`
}

type xmlWorksheet struct {
	Table xmlTable `xml:"Table"`
}

type xmlTable struct {
	Rows []xmlRow `xml:"Row"`
}

type xmlRow struct {
	Cells []xmlCell `xml:"Cell"`
}

type xmlCell struct {
	Data xmlData `xml:"Data"`
}

type xmlData struct {
	Value string `xml:",chardata"`
}

func parseXMLRows(r io.Reader) ([][]string, error) {
	var wb xmlWorkbook
	if err := xml.NewDecoder(r).Decode(&wb); err != nil {
		return nil, fmt.Errorf("decoding SpreadsheetML: %w", err)
	}
	if len(wb.Sheets) == 0 {
		return nil, fmt.Errorf("no worksheets found")
	}

	var rows [][]string
	for _, row := range wb.Sheets[0].Table.Rows {
		var cells []string
		for _, cell := range row.Cells {
			cells = append(cells, strings.TrimSpace(cell.Data.Value))
		}
		rows = append(rows, cells)
	}
	return rows, nil
}

const (
	adminHQExcelURL = "https://ca.countingopinions.com/pireports/download_xls.php?f=/temp/published/published_report_75857_1000153.xls"
	branchExcelURL  = "https://ca.countingopinions.com/pireports/download_xls.php?f=/temp/published/published_report_1000396_1000167.xls"
)

// ParseAdminHQExcel downloads and parses the admin HQ Excel XML export.
func ParseAdminHQExcel() ([]LibrarySystem, error) {
	body, err := fetchHTML(adminHQExcelURL)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseAdminHQExcel(body)
}

func parseAdminHQExcel(r io.Reader) ([]LibrarySystem, error) {
	rows, err := parseXMLRows(r)
	if err != nil {
		return nil, fmt.Errorf("parsing admin HQ Excel: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("admin HQ Excel has no data rows")
	}

	colMap := buildColumnMap(rows[0], map[string][]string{
		"name":     {"location"},
		"address":  {"street address"},
		"city":     {"city"},
		"zip":      {"1.12 zip"},
		"phone":    {"phone number - administration", "public phone"},
		"email":    {"email", "contact"},
		"website":  {"web address", "web site"},
		"dirFirst": {"first name"},
		"dirLast":  {"last name"},
	})

	now := time.Now().UTC().Format(time.DateOnly)
	var systems []LibrarySystem
	seen := make(map[string]bool)

	for _, row := range rows[1:] {
		col := func(key string) string {
			i, ok := colMap[key]
			if !ok || i >= len(row) {
				return ""
			}
			return cleanText(row[i])
		}

		name := col("name")
		if name == "" {
			continue
		}
		id := Slugify(name)
		if seen[id] {
			continue
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
	}

	return systems, nil
}

// ParseBranchesExcel downloads and parses the branch/outlet Excel XML export.
func ParseBranchesExcel() ([]Branch, error) {
	body, err := fetchHTML(branchExcelURL)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseBranchesExcel(body)
}

func parseBranchesExcel(r io.Reader) ([]Branch, error) {
	rows, err := parseXMLRows(r)
	if err != nil {
		return nil, fmt.Errorf("parsing branches Excel: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("branches Excel has no data rows")
	}

	colMap := buildColumnMap(rows[0], map[string][]string{
		"location":    {"location"},
		"name":        {"10.5 name"},
		"address":     {"10.6 street address", "street address"},
		"city":        {"10.7 city"},
		"zip":         {"10.8 zip"},
		"county":      {"10.14 county", "county"},
		"phone":       {"10.15 phone", "phone"},
		"lat":         {"latitude"},
		"lng":         {"longitude"},
		"outletType":  {"10.16 outlet", "outlet type"},
		"mainLibName": {"main library name"},
		"status":      {"status"},
	})

	now := time.Now().UTC().Format(time.DateOnly)
	var branches []Branch
	seen := make(map[string]bool)

	for _, row := range rows[1:] {
		col := func(key string) string {
			i, ok := colMap[key]
			if !ok || i >= len(row) {
				return ""
			}
			return cleanText(row[i])
		}

		name := col("name")
		if name == "" {
			name = col("location")
		}
		if name == "" {
			continue
		}

		mainLibName := col("mainLibName")
		systemID := Slugify(mainLibName)

		id := Slugify(name)
		key := systemID + "/" + id
		if seen[key] {
			continue
		}
		seen[key] = true

		lat, _ := strconv.ParseFloat(col("lat"), 64)
		lng, _ := strconv.ParseFloat(col("lng"), 64)

		b := Branch{
			ID:          id,
			SystemID:    systemID,
			Name:        name,
			Address:     col("address"),
			City:        col("city"),
			ZipCode:     col("zip"),
			County:      cleanCounty(col("county")),
			Phone:       col("phone"),
			Lat:         lat,
			Lng:         lng,
			OutletType:  parseOutletType(col("outletType")),
			Status:      "open",
			LastUpdated: now,
		}
		branches = append(branches, b)
	}

	return branches, nil
}

// buildColumnMap matches header cells to logical field names.
// Each field has one or more search strings; the first header cell
// that contains a search string (case-insensitive) wins.
func buildColumnMap(headers []string, fields map[string][]string) map[string]int {
	m := make(map[string]int)
	for i, h := range headers {
		lower := strings.ToLower(h)
		for field, needles := range fields {
			if _, ok := m[field]; ok {
				continue
			}
			for _, needle := range needles {
				if strings.Contains(lower, needle) {
					m[field] = i
					break
				}
			}
		}
	}
	return m
}
