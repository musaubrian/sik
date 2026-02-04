package utils

import (
	"os"
	"path"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

func TokenizeContent(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func StemMult(words []string) ([]string, error) {
	result := []string{}

	for _, word := range words {
		stemmed, err := Stemm(word)
		if err != nil {
			return result, err
		}

		result = append(result, stemmed)
	}
	return result, nil
}

func Stemm(word string) (string, error) {
	return snowball.Stem(word, "english", false)
}

func GetSikBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := path.Join(home, ".cache", "sik")

	cleanupOldIndex(home)

	return base, nil
}

func cleanupOldIndex(home string) {
	oldBase := path.Join(home, ".sik")
	fileInfo, err := os.Stat(oldBase)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		_ = os.RemoveAll(oldBase)
	}

}

func GetIndexLocation(base string) string {
	return path.Join(base, "index.sik")
}

func Ignore(dir string) bool {
	ignoreList := []string{".git", "node_modules", ".venv"}
	return slices.Contains(ignoreList, dir)
}
