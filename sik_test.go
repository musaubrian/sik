package main

import "testing"

func TestReadMarkdown(t *testing.T) {
	testPaths := []struct {
		path string
		want bool
	}{
		{path: "test_dir/basics.md", want: true},
		{path: "test_dir/project.md", want: true},
	}

	res, err := readMarkdown("./test_dir")
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
	res, _ := readMarkdown("./test_dir")

	_, err := createIndex(res)
	if err != nil {
		t.Fatalf("Expected <nil> got <%s>", err)
	}
}

func TestIndex(t *testing.T) {
	res, _ := readMarkdown("./test_dir")

	index, _ := createIndex(res)

	singleOccurrenceWords := []struct {
		word string
		file string
	}{
		{"basic", "test_dir/basics.md"},
		{"brief", "test_dir/project.md"},
	}

	for _, tt := range singleOccurrenceWords {
		if fileMeta, exists := index[tt.word]; !exists {
			t.Errorf("Expected word '%s' to be in the index", tt.word)
		} else {
			if count, fileExists := fileMeta[tt.file]; !fileExists || count != 1 {
				t.Errorf("Expected word '%s' to appear once in '%s', got %d", tt.word, tt.file, count)
			}
		}
	}

	multiOccurenceWords := []struct {
		word string
		file string
	}{
		{"stuff", "test_dir/basics.md"},
		{"featur", "test_dir/project.md"},
	}

	for _, tt := range multiOccurenceWords {
		if fileMeta, exists := index[tt.word]; !exists {
			t.Errorf("Expected word '%s' to be in the index", tt.word)
		} else {
			if count, fileExists := fileMeta[tt.file]; !fileExists || count <= 1 {
				t.Errorf("Expected word '%s' to appear once in '%s', got %d", tt.word, tt.file, count)
			}
		}
	}
}

func TestSearch(t *testing.T) {
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
