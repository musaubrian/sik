package main

import (
	"os"
	"path"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

func tokenizeContent(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func stemm(word string) (string, error) {
	return snowball.Stem(word, "english", true)
}

func getSikBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := path.Join(home, ".sik")

	return base, nil
}

func getIndexLocation(base string) string {
	return path.Join(base, "index.sik")
}

func ignore(dir string) bool {
	ignoreList := []string{".git", "node_modules", ".venv"}
	return slices.Contains(ignoreList, dir)
}
