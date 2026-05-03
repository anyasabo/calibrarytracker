package scraper

import "testing"

func TestFormatPhone(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"9162642920", "(916) 264-2920"},
		{"19162642920", "(916) 264-2920"},
		{"(916) 264-2920", "(916) 264-2920"},
		{"916-264-2920", "(916) 264-2920"},
		{"", ""},
		{"12345", "12345"},
	}
	for _, tt := range tests {
		got := formatPhone(tt.input)
		if got != tt.want {
			t.Errorf("formatPhone(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseBranchStatus(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"no change", "open"},
		{"open", "open"},
		{"", "open"},
		{"closed", "closed"},
		{"permanently closed", "closed"},
		{"temporarily closed", "temporarily-closed"},
		{"Temporarily Closed", "temporarily-closed"},
	}
	for _, tt := range tests {
		got := parseBranchStatus(tt.input)
		if got != tt.want {
			t.Errorf("parseBranchStatus(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidateCACoords(t *testing.T) {
	tests := []struct {
		name     string
		lat, lng float64
		wantLat  float64
		wantLng  float64
	}{
		{"Sacramento", 38.58, -121.49, 38.58, -121.49},
		{"zero", 0, 0, 0, 0},
		{"outside CA lat", 50.0, -121.0, 0, 0},
		{"outside CA lng", 37.0, -80.0, 0, 0},
		{"negative lat", -38.0, -121.0, 0, 0},
	}
	for _, tt := range tests {
		gotLat, gotLng := validateCACoords(tt.lat, tt.lng, tt.name)
		if gotLat != tt.wantLat || gotLng != tt.wantLng {
			t.Errorf("validateCACoords(%f, %f, %q) = (%f, %f), want (%f, %f)",
				tt.lat, tt.lng, tt.name, gotLat, gotLng, tt.wantLat, tt.wantLng)
		}
	}
}
