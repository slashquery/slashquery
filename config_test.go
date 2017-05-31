package slashquery

import "testing"

func TestConfigNew(t *testing.T) {
	var testTable = []struct {
		yml      string
		expected bool
	}{
		{"testdata/non-existent.yml", false},
		{"testdata/slashquery.yml", true},
	}
	for _, tt := range testTable {
		_, err := New(tt.yml)
		if !tt.expected {
			if err == nil {
				t.Fatal(err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
