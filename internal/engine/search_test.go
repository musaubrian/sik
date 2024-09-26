package engine

import (
	"testing"

	"github.com/musaubrian/sik/internal/core"
	"github.com/musaubrian/sik/internal/utils"
)

func TestSimpleSearch(t *testing.T) {
	res, _ := core.ReadMarkdown("../test_src")
	in, _ := core.CreateIndex(res)

	words := []string{"duplicate", "features", "description"}

	s := New(in)

	for _, word := range words {
		searchRes, err := s.Search(word)
		if err != nil {
			t.Fatalf("Expected <nil>, Got <%s>", err)
		}
		if len(searchRes) == 0 {
			t.Fatalf("Expected 1 result, got <%d>", len(searchRes))
		}
	}
}

func TestPhraseSearch(t *testing.T) {
	res, _ := core.ReadMarkdown("../test_src")
	ind, _ := core.CreateIndex(res)

	s := New(ind)

	testCases := []struct {
		phrase   string
		expected []string
	}{
		{"description of", []string{"../test_src/basics.md", "../test_src/project.md"}},
		{"Basic stuff", []string{"../test_src/basics.md"}},
		{"more stuff", []string{"../test_src/basics.md"}},
		{"of markdown", []string{"../test_src/basics.md"}},
	}
	for _, tc := range testCases {
		searchRes, err := s.phraseSearch(utils.TokenizeContent(tc.phrase))
		if err != nil {
			t.Errorf("Expect <nil>, Got: %v", err)
		}

		if !unorderedEqual(searchRes, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.phrase, searchRes, tc.expected)
		}
	}
}

func TestProximitySearch(t *testing.T) {

	res, _ := core.ReadMarkdown("../test_src")
	ind, _ := core.CreateIndex(res)

	s := New(ind)
	testCases := []struct {
		query    string
		expected []string
	}{
		{"brief project", []string{"../test_src/project.md"}},
		{"description duplicate", []string{"../test_src/basics.md", "../test_src/project.md"}},
		{"description markdown", []string{"../test_src/basics.md"}},
	}
	for _, tc := range testCases {
		searchRes, err := s.proximitySearch(utils.TokenizeContent(tc.query))
		if err != nil {
			t.Errorf("Expect <nil>, Got: %v", err)
		}

		if !unorderedEqual(searchRes, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.query, searchRes, tc.expected)
		}
	}
}

// func equal(a, b []string) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i := range a {
// 		if a[i] != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

func unorderedEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int, len(a))
	for _, s := range a {
		m[s]++
	}
	for _, s := range b {
		if m[s] == 0 {
			return false
		}
		m[s]--
	}
	return true
}
