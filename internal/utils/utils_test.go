package utils

import "testing"

func TestGetSikBase(t *testing.T) {
	_, err := GetSikBase()

	if err != nil {
		t.Fatalf("Expected <nil> got <%v>", err)
	}
}

func TestTokenizing(t *testing.T) {
	multiWords := []struct {
		text        string
		expectedLen int
	}{
		{"hello there, how are you bruv", 6},
		{"no pnctuation", 2},
		{"So, this; will remove : punctuation ?", 5},
		{"This + 90 equal 10294 .", 4},
		{"one", 1},
		{": , . > | *", 0},
	}

	for _, tt := range multiWords {
		multiToken := TokenizeContent(tt.text)
		if len(multiToken) != tt.expectedLen {
			t.Errorf("Expected tokenized content to be %d, Got: %d", tt.expectedLen, len(multiToken))
		}
	}
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
		if got := Ignore(tt.dir); got != tt.want {
			t.Errorf("ignore(%v) -> %v; Expected %v", tt.dir, got, tt.want)
		}
	}
}
