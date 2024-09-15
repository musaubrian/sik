package main

import "testing"

func TestGetSikBase(t *testing.T) {
	_, err := getSikBase()

	if err != nil {
		t.Fatalf("Expected <nil> got <%v>", err)
	}
}

func TestTokenizing(t *testing.T) {

}

func TestIgnore(t *testing.T) {
	cases := []struct {
		dir  string
		want bool
	}{
		{dir: ".git", want: true},
		{dir: "src", want: false},
		{dir: "node_modules", want: true},
		{dir: "bin", want: false},
		{dir: ".venv", want: true},
		{dir: "custom_dir", want: false},
	}

	for _, tt := range cases {
		if got := ignore(tt.dir); got != tt.want {
			t.Errorf("ignore(%v) -> %v; Expected %v", tt.dir, got, tt.want)
		}
	}
}
