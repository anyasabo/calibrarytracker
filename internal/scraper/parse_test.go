package scraper

import (
	"strings"
	"testing"
)

const sampleAdminHQHTML = `<html><body>
<table>
<thead><tr>
<th>Location (186)</th><th>1.10 Street Address</th><th>1.11 City</th><th>1.12 Zip</th>
<th>1.13 Zip +4</th><th>1.18 Public Phone Number - Administration</th>
<th>1.19 Reference Phone Number</th><th>1.20 TDD for Deaf</th>
<th>1.28 Library Public Email address or "contact us" URL</th>
<th>1.29 Library's Web Address</th><th>1.6 Director First Name</th>
<th>1.8 Director Last Name</th><th>1.9 Director Title</th>
</tr></thead>
<tbody>
<tr>
<td>SACRAMENTO PUBLIC LIBRARY</td><td>828 I ST.</td><td>SACRAMENTO</td><td>95814</td>
<td>2508</td><td>(916) 264-2920</td><td>(916) 264-2920</td><td>(916) 264-2855</td>
<td>contact@saclibrary.org</td><td>http://www.saclibrary.org</td>
<td>Peter</td><td>Coyl</td><td>Library Director &amp; CEO</td>
</tr>
<tr>
<td>SAN FRANCISCO PUBLIC LIBRARY</td><td>100 LARKIN ST</td><td>SAN FRANCISCO</td><td>94102</td>
<td>4705</td><td>(415) 557-4400</td><td>(415) 557-4400</td><td>(415) 557-4433</td>
<td>info@sfpl.org</td><td>http://sfpl.org</td>
<td>Michael</td><td>Lambert</td><td>City Librarian</td>
</tr>
</tbody>
</table>
</body></html>`

func TestParseAdminHQHTML(t *testing.T) {
	systems, err := parseAdminHQHTML(strings.NewReader(sampleAdminHQHTML))
	if err != nil {
		t.Fatalf("parseAdminHQHTML() error: %v", err)
	}
	if len(systems) != 2 {
		t.Fatalf("expected 2 systems, got %d", len(systems))
	}

	sac := systems[0]
	if sac.Name != "SACRAMENTO PUBLIC LIBRARY" {
		t.Errorf("name = %q, want SACRAMENTO PUBLIC LIBRARY", sac.Name)
	}
	if sac.City != "SACRAMENTO" {
		t.Errorf("city = %q, want SACRAMENTO", sac.City)
	}
	if sac.Address != "828 I ST." {
		t.Errorf("address = %q, want 828 I ST.", sac.Address)
	}
	if sac.Phone != "(916) 264-2920" {
		t.Errorf("phone = %q, want (916) 264-2920", sac.Phone)
	}
	if sac.Website != "http://www.saclibrary.org" {
		t.Errorf("website = %q, want http://www.saclibrary.org", sac.Website)
	}
	if sac.Email != "contact@saclibrary.org" {
		t.Errorf("email = %q, want contact@saclibrary.org", sac.Email)
	}
	if sac.DirectorName != "Peter Coyl" {
		t.Errorf("directorName = %q, want Peter Coyl", sac.DirectorName)
	}
	if sac.ID != "sacramento-public-library" {
		t.Errorf("id = %q, want sacramento-public-library", sac.ID)
	}

	sf := systems[1]
	if sf.Name != "SAN FRANCISCO PUBLIC LIBRARY" {
		t.Errorf("sf name = %q, want SAN FRANCISCO PUBLIC LIBRARY", sf.Name)
	}
}

const sampleBranchHTML = `<html><body>
<table>
<thead><tr>
<th>Location</th><th>10.5 Name</th><th>10.6 Street Address</th><th>10.7 City</th>
<th>10.8 Zip Code</th><th>10.9 Zip+4 Code</th><th>10.10 Mailing Street Address</th>
<th>10.11 Mailing City</th><th>10.12 Mailing Zip Code</th><th>10.13 Mailing Zip +4</th>
<th>10.14 County</th><th>10.15 Phone</th><th>Latitude</th><th>Longitude</th>
<th>10.16 Outlet Type Code</th><th>Main Library Name</th>
<th>Status of Outlet record</th><th>Law Library</th>
</tr></thead>
<tbody>
<tr>
<td>CENTRAL LIBRARY</td><td>CENTRAL LIBRARY</td><td>828 I ST.</td><td>SACRAMENTO</td>
<td>95814</td><td>2508</td><td>828 I ST.</td><td>SACRAMENTO</td>
<td>95814</td><td>2508</td><td>34Sacramento</td><td>9162642920</td>
<td>38.58157</td><td>-121.49440</td><td>2Central</td><td>SACRAMENTO PUBLIC LIBRARY</td>
<td>0no change</td><td>1No</td>
</tr>
<tr>
<td>ARDEN-DIMICK LIBRARY</td><td>ARDEN-DIMICK LIBRARY</td><td>891 WATT AVE.</td><td>SACRAMENTO</td>
<td>95864</td><td>4621</td><td>891 WATT AVE.</td><td>SACRAMENTO</td>
<td>95864</td><td>4621</td><td>34Sacramento</td><td>9162647100</td>
<td>38.59777</td><td>-121.39706</td><td>0Branch</td><td>SACRAMENTO PUBLIC LIBRARY</td>
<td>0no change</td><td>1No</td>
</tr>
</tbody>
</table>
</body></html>`

func TestParseBranchesHTML(t *testing.T) {
	branches, err := parseBranchesHTML(strings.NewReader(sampleBranchHTML))
	if err != nil {
		t.Fatalf("parseBranchesHTML() error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	central := branches[0]
	if central.Name != "CENTRAL LIBRARY" {
		t.Errorf("name = %q, want CENTRAL LIBRARY", central.Name)
	}
	if central.SystemID != "sacramento-public-library" {
		t.Errorf("systemId = %q, want sacramento-public-library", central.SystemID)
	}
	if central.City != "SACRAMENTO" {
		t.Errorf("city = %q, want SACRAMENTO", central.City)
	}
	if central.Address != "828 I ST." {
		t.Errorf("address = %q, want 828 I ST.", central.Address)
	}
	if central.County != "Sacramento" {
		t.Errorf("county = %q, want Sacramento", central.County)
	}
	if central.OutletType != "central" {
		t.Errorf("outletType = %q, want central", central.OutletType)
	}
	if central.Lat == 0 || central.Lng == 0 {
		t.Errorf("lat/lng should not be zero: %f, %f", central.Lat, central.Lng)
	}

	arden := branches[1]
	if arden.OutletType != "branch" {
		t.Errorf("outletType = %q, want branch", arden.OutletType)
	}
}
