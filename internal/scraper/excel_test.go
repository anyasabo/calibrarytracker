package scraper

import (
	"strings"
	"testing"
)

const sampleAdminHQExcelXML = `<?xml version="1.0" encoding="UTF-8"?>
<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet">
<Worksheet><Table>
<Row>
<Cell><Data ss:Type="String">Location</Data></Cell>
<Cell><Data ss:Type="String">1.10 Street Address</Data></Cell>
<Cell><Data ss:Type="String">1.11 City</Data></Cell>
<Cell><Data ss:Type="String">1.12 Zip</Data></Cell>
<Cell><Data ss:Type="String">1.13 Zip +4</Data></Cell>
<Cell><Data ss:Type="String">1.18 Public Phone Number - Administration</Data></Cell>
<Cell><Data ss:Type="String">1.19 Reference Phone Number</Data></Cell>
<Cell><Data ss:Type="String">1.20 TDD for Deaf</Data></Cell>
<Cell><Data ss:Type="String">1.28 Library Public Email address</Data></Cell>
<Cell><Data ss:Type="String">1.29 Library's Web Address</Data></Cell>
<Cell><Data ss:Type="String">1.6 Director First Name</Data></Cell>
<Cell><Data ss:Type="String">1.8 Director Last Name</Data></Cell>
<Cell><Data ss:Type="String">1.9 Director Title</Data></Cell>
</Row>
<Row>
<Cell><Data ss:Type="String">SACRAMENTO PUBLIC LIBRARY</Data></Cell>
<Cell><Data ss:Type="String">828 I ST.</Data></Cell>
<Cell><Data ss:Type="String">SACRAMENTO</Data></Cell>
<Cell><Data ss:Type="String">95814</Data></Cell>
<Cell><Data ss:Type="String">2508</Data></Cell>
<Cell><Data ss:Type="String">(916) 264-2920</Data></Cell>
<Cell><Data ss:Type="String">(916) 264-2920</Data></Cell>
<Cell><Data ss:Type="String">(916) 264-2855</Data></Cell>
<Cell><Data ss:Type="String">contact@saclibrary.org</Data></Cell>
<Cell><Data ss:Type="String">http://www.saclibrary.org</Data></Cell>
<Cell><Data ss:Type="String">Peter</Data></Cell>
<Cell><Data ss:Type="String">Coyl</Data></Cell>
<Cell><Data ss:Type="String">Library Director &amp; CEO</Data></Cell>
</Row>
</Table></Worksheet>
</Workbook>`

func TestParseAdminHQExcel(t *testing.T) {
	systems, err := parseAdminHQExcel(strings.NewReader(sampleAdminHQExcelXML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(systems) != 1 {
		t.Fatalf("got %d systems, want 1", len(systems))
	}

	s := systems[0]
	checks := []struct{ field, got, want string }{
		{"ID", s.ID, "sacramento-public-library"},
		{"Name", s.Name, "SACRAMENTO PUBLIC LIBRARY"},
		{"Address", s.Address, "828 I ST."},
		{"City", s.City, "SACRAMENTO"},
		{"Phone", s.Phone, "(916) 264-2920"},
		{"Email", s.Email, "contact@saclibrary.org"},
		{"Website", s.Website, "http://www.saclibrary.org"},
		{"DirectorName", s.DirectorName, "Peter Coyl"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}
}

const sampleBranchExcelXML = `<?xml version="1.0" encoding="UTF-8"?>
<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet">
<Worksheet><Table>
<Row>
<Cell><Data ss:Type="String">Location</Data></Cell>
<Cell><Data ss:Type="String">10.5 Name</Data></Cell>
<Cell><Data ss:Type="String">10.6 Street Address</Data></Cell>
<Cell><Data ss:Type="String">10.7 City</Data></Cell>
<Cell><Data ss:Type="String">10.8 Zip Code</Data></Cell>
<Cell><Data ss:Type="String">10.9 Zip+4 Code</Data></Cell>
<Cell><Data ss:Type="String">10.10 Mailing Street Address</Data></Cell>
<Cell><Data ss:Type="String">10.11 Mailing City</Data></Cell>
<Cell><Data ss:Type="String">10.12 Mailing Zip Code</Data></Cell>
<Cell><Data ss:Type="String">10.13 Mailing Zip +4</Data></Cell>
<Cell><Data ss:Type="String">10.14 County</Data></Cell>
<Cell><Data ss:Type="String">10.15 Phone</Data></Cell>
<Cell><Data ss:Type="String">Latitude</Data></Cell>
<Cell><Data ss:Type="String">Longitude</Data></Cell>
<Cell><Data ss:Type="String">10.16 Outlet Type Code</Data></Cell>
<Cell><Data ss:Type="String">Main Library Name</Data></Cell>
<Cell><Data ss:Type="String">Status</Data></Cell>
<Cell><Data ss:Type="String">Law Library</Data></Cell>
</Row>
<Row>
<Cell><Data ss:Type="String">ARCADE LIBRARY</Data></Cell>
<Cell><Data ss:Type="String">ARCADE LIBRARY</Data></Cell>
<Cell><Data ss:Type="String">2443 MARCONI AVE.</Data></Cell>
<Cell><Data ss:Type="String">SACRAMENTO</Data></Cell>
<Cell><Data ss:Type="String">95821</Data></Cell>
<Cell><Data ss:Type="String">4030</Data></Cell>
<Cell><Data ss:Type="String">2443 MARCONI AVE.</Data></Cell>
<Cell><Data ss:Type="String">SACRAMENTO</Data></Cell>
<Cell><Data ss:Type="String">95821</Data></Cell>
<Cell><Data ss:Type="String">4030</Data></Cell>
<Cell><Data ss:Type="String">Sacramento</Data></Cell>
<Cell><Data ss:Type="String">9162642920</Data></Cell>
<Cell><Data ss:Type="String">38.61861</Data></Cell>
<Cell><Data ss:Type="String">-121.40409</Data></Cell>
<Cell><Data ss:Type="String">Branch</Data></Cell>
<Cell><Data ss:Type="String">SACRAMENTO PUBLIC LIBRARY</Data></Cell>
<Cell><Data ss:Type="String">no change</Data></Cell>
<Cell><Data ss:Type="String">No</Data></Cell>
</Row>
</Table></Worksheet>
</Workbook>`

func TestParseBranchesExcel(t *testing.T) {
	branches, err := parseBranchesExcel(strings.NewReader(sampleBranchExcelXML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("got %d branches, want 1", len(branches))
	}

	b := branches[0]
	checks := []struct{ field, got, want string }{
		{"ID", b.ID, "arcade-library"},
		{"SystemID", b.SystemID, "sacramento-public-library"},
		{"Name", b.Name, "ARCADE LIBRARY"},
		{"Address", b.Address, "2443 MARCONI AVE."},
		{"City", b.City, "SACRAMENTO"},
		{"ZipCode", b.ZipCode, "95821"},
		{"County", b.County, "Sacramento"},
		{"Phone", b.Phone, "(916) 264-2920"},
		{"OutletType", b.OutletType, "branch"},
		{"Status", b.Status, "open"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.field, c.got, c.want)
		}
	}

	if b.Lat != 38.61861 {
		t.Errorf("Lat = %f, want 38.61861", b.Lat)
	}
	if b.Lng != -121.40409 {
		t.Errorf("Lng = %f, want -121.40409", b.Lng)
	}
}
