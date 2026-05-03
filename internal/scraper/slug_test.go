package scraper

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Sacramento Public Library", "sacramento-public-library"},
		{"A. K. SMILEY PUBLIC LIBRARY", "a-k-smiley-public-library"},
		{"  Leading and Trailing Spaces  ", "leading-and-trailing-spaces"},
		{"LA COUNTY LIBRARY", "la-county-library"},
		{"San José Public Library", "san-jos-public-library"},
		{"81ST AVENUE BRANCH LIBRARY", "81st-avenue-branch-library"},
		{"", ""},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
