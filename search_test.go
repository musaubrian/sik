package main

import "testing"

func TestSimpleSearch(t *testing.T) {
	res, _ := readMarkdown("./test_dir")
	in, _ := createIndex(res)

	words := []string{"duplicate", "features", "description"}

	for _, word := range words {
		searchRes, err := search(word, in)
		if err != nil {
			t.Fatalf("Expected <nil>, Got <%s>", err)
		}
		if len(searchRes) == 0 {
			t.Fatalf("Expected 1 result, got <%d>", len(searchRes))
		}
	}
}

func TestPhraseSearch(t *testing.T) {
	res, _ := readMarkdown("./test_dir")

	ind, _ := createIndex(res)

	testCases := []struct {
		phrase   string
		expected []string
	}{
		{"description of", []string{"test_dir/basics.md", "test_dir/project.md"}},
		{"Basic stuff", []string{"test_dir/basics.md"}},
		{"more stuff", []string{"test_dir/basics.md"}},
		{"of markdown", []string{"test_dir/basics.md"}},
	}
	for _, tc := range testCases {
		searchRes, err := phraseSearch(tokenizeContent(tc.phrase), ind)
		if err != nil {
			t.Errorf("Expect <nil>, Got: %v", err)
		}

		if !equal(searchRes, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.phrase, searchRes, tc.expected)
		}
	}
}

func TestProximitySearch(t *testing.T) {

	res, _ := readMarkdown("./test_dir")

	ind, _ := createIndex(res)
	testCases := []struct {
		query    string
		expected []string
	}{
		{"brief project", []string{"test_dir/project.md"}},
		{"description duplicate", []string{"test_dir/basics.md", "test_dir/project.md"}},
		{"description markdown", []string{"test_dir/basics.md"}},
	}
	for _, tc := range testCases {
		searchRes, err := proximitySearch(tokenizeContent(tc.query), ind)
		if err != nil {
			t.Errorf("Expect <nil>, Got: %v", err)
		}

		if !equal(searchRes, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.query, searchRes, tc.expected)
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
