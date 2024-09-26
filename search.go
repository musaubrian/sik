package main

import (
	"fmt"
	"log/slog"
	"math"
	"strings"
)

const MAX_DISTANCE = 10 // How far in either direction to check for words

func search(query string, index Index) ([]string, error) {
	// var error error

	tContent := tokenizeContent(query)
	if len(tContent) == 1 {
		return simpleSearch(tContent[0], index)

	}

	results, err := phraseSearch(tContent, index)
	if err != nil {
		slog.Error(err.Error())
		// error = fmt.Errorf("[phrase_search]: %v", err)
	}

	if len(results) == 0 {
		return proximitySearch(tContent, index)
	}

	return results, nil
}

func simpleSearch(query string, index Index) ([]string, error) {
	stemmed, err := stemm(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to stemm word: %v", err)
	}

	results := []string{}
	for k, meta := range index {
		if strings.Contains(k, stemmed) {
			for path := range meta {
				results = append(results, path)
			}
		}
	}

	return results, nil
}

func phraseSearch(phrase []string, index Index) ([]string, error) {
	res := []string{}
	set := []FileMeta{}

	stemmedPhrase, err := stemMult(phrase)
	if err != nil {
		return res, fmt.Errorf("Could not stemm words: %v", err)
	}

	for _, word := range stemmedPhrase {
		if meta, found := index[word]; found {
			set = append(set, meta)
		} else {
			return res, nil
		}
	}

	commonDocs := intersectAll(set)

	for doc := range commonDocs {
		if wordsAppearInSequence(doc, stemmedPhrase, index) {
			res = append(res, doc)
		}
	}

	return res, nil
}

func proximitySearch(query []string, index Index) ([]string, error) {
	finalResult := []string{}
	resultsSet := []FileMeta{}

	stemmedQuery, err := stemMult(query)
	if err != nil {
		return finalResult, fmt.Errorf("Could not stemm words: %v", err)
	}

	for _, word := range stemmedQuery {

		if fileMeta, ok := index[word]; ok {
			resultsSet = append(resultsSet, fileMeta)
		} else {
			return finalResult, nil
		}
	}

	commonDocs := intersectAll(resultsSet)

	for doc := range commonDocs {
		if wordsInProximity(doc, stemmedQuery, index) {
			finalResult = append(finalResult, doc)
		}

	}
	return finalResult, nil
}

func intersectAll(resultSets []FileMeta) map[string]struct{} {
	commonDocs := make(map[string]struct{})
	if len(resultSets) == 0 {
		return commonDocs
	}

	for doc := range resultSets[0] {
		commonDocs[doc] = struct{}{}
	}

	for _, resultSet := range resultSets[1:] {
		for doc := range commonDocs {
			if _, found := resultSet[doc]; !found {
				delete(commonDocs, doc) // Remove if not in current set
			}
		}
	}
	return commonDocs
}

func wordsAppearInSequence(doc string, queryWords []string, index Index) bool {
	firstWord := queryWords[0]
	firstWordPositions := index[firstWord][doc]

	// For each occurrence of the first word, check if the rest of the words follow in sequence
	for _, startPos := range firstWordPositions {
		matched := true
		for i := 1; i < len(queryWords); i++ {
			nextWord := queryWords[i]
			nextWordPositions := index[nextWord][doc]
			expectedPos := startPos + i

			// Check if any of the positions match the expected position
			if !containsPosition(nextWordPositions, expectedPos) {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}

	return false
}

func wordsInProximity(doc string, words []string, index Index) bool {
	wordPositions := make([][]int, len(words))
	for i, word := range words {
		if pos, ok := index[word][doc]; ok {
			wordPositions[i] = pos
		} else {
			return false
		}
	}

	for i := 0; i < len(wordPositions)-1; i++ {
		for _, pos1 := range wordPositions[i] {
			for _, pos2 := range wordPositions[i+1] {
				if math.Abs(float64(pos1)-float64(pos2)) <= float64(MAX_DISTANCE) {
					return true
				}
			}
		}
	}

	return false
}

func containsPosition(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}
