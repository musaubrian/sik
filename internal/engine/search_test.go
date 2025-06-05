package engine

import (
	"testing"

	"github.com/musaubrian/sik/internal/core"
	"github.com/musaubrian/sik/internal/utils"
)

func TestSimpleSearch(t *testing.T) {
	res, _ := core.ReadMarkdown("../test_src")
	in, _ := core.CreateIndex(res)

	words := []string{"duplicate", "features", "description", "brief"}

	s := New(in)

	for _, word := range words {
		searchRes, err := s.simpleSearch(word)
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
		{"description of", []string{"../test_src/project.md", "../test_src/basics.md"}},
		{"Basic stuff", []string{"../test_src/basics.md"}},
		{"more stuff", []string{"../test_src/basics.md"}},
		{"of markdown", []string{"../test_src/basics.md"}},
	}
	for _, tc := range testCases {
		docsRes, err := s.phraseSearch(utils.TokenizeContent(tc.phrase))
		if err != nil {
			t.Errorf("Expect <nil>, Got: %v", err)
		}
		normDocs := normalizeDocRes(docsRes)

		if !orderedEqual(normDocs, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.phrase, normDocs, tc.expected)
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
		{"description duplicate", []string{"../test_src/project.md", "../test_src/basics.md"}},
		{"description markdown", []string{"../test_src/basics.md"}},
	}
	for _, tc := range testCases {
		docsRes, err := s.proximitySearch(utils.TokenizeContent(tc.query))
		if err != nil {
			t.Errorf("Expect <nil>, Got: %v", err)
		}

		normDocs := normalizeDocRes(docsRes)

		if !orderedEqual(normDocs, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.query, normDocs, tc.expected)
		}
	}
}

func TestGeneralSearch(t *testing.T) {
	res, _ := core.ReadMarkdown("../test_src")
	ind, _ := core.CreateIndex(res)

	s := New(ind)
	testCases := []struct {
		query    string
		expected []string
	}{
		{"brief", []string{"../test_src/basics.md", "../test_src/project.md"}},
		{"brief project", []string{"../test_src/project.md"}},
		{"description duplicate", []string{"../test_src/project.md", "../test_src/basics.md"}},
		{"description markdown", []string{"../test_src/basics.md"}},
	}
	for _, tc := range testCases {
		searchRes, err := s.Search(tc.query)
		if err != nil {
			t.Errorf("Expect <nil> got %v", err)
		}

		if len(searchRes) < 0 {
			t.Errorf("Expected non empty result")
		}

		if !orderedEqual(searchRes, tc.expected) {
			t.Errorf("For: <%s>; got %v, want %v", tc.query, searchRes, tc.expected)
		}
	}
}

func orderedEqual(a, b []string) bool {
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

func normalizeDocRes(docs []DocRes) []string {
	vals := []string{}
	for _, doc := range docs {
		vals = append(vals, doc.Path)
	}

	return vals
}
