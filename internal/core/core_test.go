package core

import (
	"testing"

	"github.com/musaubrian/sik/internal/utils"
)

func TestReadMarkdown(t *testing.T) {
	testPaths := []struct {
		path string
		want bool
	}{
		{path: "../test_src/basics.md", want: true},
		{path: "../test_src/project.md", want: true},
	}

	res, err := ReadMarkdown("../test_src")
	if err != nil {
		t.Fatalf("Expected <nil>, got <%v>", err)
	}

	for _, tt := range testPaths {
		if key, ok := res[tt.path]; !ok {
			t.Fatalf("Expected <%s>, Got <%s>", tt.path, key)
		}
	}

}

func TestCreateIndex(t *testing.T) {
	res, _ := ReadMarkdown("./../test_src")

	_, err := CreateIndex(res)
	if err != nil {
		t.Fatalf("Expected <nil> got <%s>", err)
	}
}

func TestIndex(t *testing.T) {
	res, _ := ReadMarkdown("../test_src")

	index, _ := CreateIndex(res)

	singleOccurrenceWords := []struct {
		word string
		file string
	}{
		{"basic", "../test_src/basics.md"},
		{"brief", "../test_src/project.md"},
	}

	for _, tt := range singleOccurrenceWords {
		stemmedWord, err := utils.Stemm(tt.word)
		if err != nil {
			Logging.Error(err.Error())
		}
		if fileMeta, exists := index[stemmedWord]; !exists {
			t.Errorf("Expected word '%s' to be in the index", tt.word)
		} else {
			if pos, exists := fileMeta[tt.file]; !exists || len(pos) != 1 {
				t.Errorf("Expected word <%s> to appear once in <%s>, Got %d", tt.word, tt.file, len(pos))
			}
		}
	}

	multiOccurenceWords := []struct {
		word string
		file string
	}{
		{"stuff", "../test_src/basics.md"},
		{"feature", "../test_src/project.md"},
	}

	for _, tt := range multiOccurenceWords {
		stemmedWord, err := utils.Stemm(tt.word)
		if err != nil {
			Logging.Error(err.Error())
		}
		if fileMeta, exists := index[stemmedWord]; !exists {
			t.Errorf("Expected word '%s' to be in the index", tt.word)
		} else {
			if pos, exists := fileMeta[tt.file]; !exists || len(pos) <= 1 {
				t.Errorf("Expected word '%s' to appear once in '%s', got %d", tt.word, tt.file, len(pos))
			}
		}
	}
}
