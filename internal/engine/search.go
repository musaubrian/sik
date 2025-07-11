package engine

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/musaubrian/sik/internal/core"
	"github.com/musaubrian/sik/internal/utils"
)

const maxProximityDistance = 10

type Engine struct {
	index                core.IndexContents
	maxProximityDistance int
}

type DocRes struct {
	Path       string
	Occurences int
}

func New(index core.IndexContents) *Engine {
	return &Engine{
		index:                index,
		maxProximityDistance: maxProximityDistance,
	}
}

func removeDuplicates(in []DocRes) []string {
	clean := []string{}

	slices.SortStableFunc(in, cmp)

	for _, v := range in {
		if !slices.Contains(clean, v.Path) {
			clean = append(clean, v.Path)
		}
	}

	return clean
}

func (se *Engine) Search(query string) ([]string, error) {
	tokens := utils.TokenizeContent(query)
	if len(tokens) == 1 {
		simpleSearchRes, err := se.simpleSearch(tokens[0])
		return removeDuplicates(simpleSearchRes), err
	}

	phraseResults, err := se.phraseSearch(tokens)
	if err != nil {
		return nil, fmt.Errorf("[phrase_search]: %v", err)
	}
	if len(phraseResults) > 0 {
		return removeDuplicates(phraseResults), nil
	}

	proximityResults, err := se.proximitySearch(tokens)
	if err != nil {
		return nil, fmt.Errorf("[proximity_search]: %v", err)
	}
	return removeDuplicates(proximityResults), nil
}

func (se *Engine) simpleSearch(query string) ([]DocRes, error) {
	stemmed, err := utils.Stemm(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to stemm word: %v", err)
	}

	results := []DocRes{}
	for k, meta := range se.index {
		if strings.Contains(k, stemmed) {
			for path, occurences := range meta {
				results = append(results, DocRes{Path: path, Occurences: len(occurences)})
			}
		}
	}

	return sortDocs(results), nil
}

func (se *Engine) phraseSearch(phrase []string) ([]DocRes, error) {
	res := []DocRes{}
	set := []core.FileMeta{}

	stemmedPhrase, err := utils.StemMult(phrase)
	if err != nil {
		return res, fmt.Errorf("Could not stemm words: %v", err)
	}

	for _, word := range stemmedPhrase {
		if meta, found := se.index[word]; found {
			set = append(set, meta)
		} else {
			return res, nil
		}
	}

	commonDocs := mergeAll(set)

	for doc, occ := range commonDocs {
		if wordsAppearInSequence(doc, stemmedPhrase, se.index) {
			res = append(res, DocRes{Path: doc, Occurences: len(occ)})
		}
	}

	return sortDocs(res), nil
}

func (se *Engine) proximitySearch(query []string) ([]DocRes, error) {
	finalResult := []DocRes{}
	resultsSet := []core.FileMeta{}

	stemmedQuery, err := utils.StemMult(query)
	if err != nil {
		return finalResult, fmt.Errorf("Could not stemm words: %v", err)
	}

	for _, word := range stemmedQuery {

		if fileMeta, ok := se.index[word]; ok {
			resultsSet = append(resultsSet, fileMeta)
		} else {
			return finalResult, nil
		}
	}

	commonDocs := mergeAll(resultsSet)

	for doc, occ := range commonDocs {
		if wordsInProximity(doc, stemmedQuery, se.index, se.maxProximityDistance) {
			finalResult = append(finalResult, DocRes{Path: doc, Occurences: len(occ)})
		}

	}
	return sortDocs(finalResult), nil
}

func mergeAll(resultSets []core.FileMeta) map[string][]int {
	commonDocs := make(map[string][]int)
	if len(resultSets) == 0 {
		return commonDocs
	}

	for _, meta := range resultSets {
		for path, occ := range meta {
			if val, foundExisting := commonDocs[path]; !foundExisting {
				commonDocs[path] = val
			}

			commonDocs[path] = append(commonDocs[path], occ...)
		}
	}

	return commonDocs
}

func wordsAppearInSequence(doc string, queryWords []string, index core.IndexContents) bool {
	firstWord := queryWords[0]
	firstWordPositions := index[firstWord][doc]

	// For each occurrence of the first word, check if the rest of the words follow in sequence
	for _, startPos := range firstWordPositions {
		matched := true
		for i := 1; i < len(queryWords); i++ {
			nextWord := queryWords[i]
			nextWordPositions := index[nextWord][doc]
			expectedPos := startPos + i

			if !slices.Contains(nextWordPositions, expectedPos) {
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

func wordsInProximity(doc string, words []string, index core.IndexContents, maxDistance int) bool {
	wordPositions := make([][]int, len(words))
	for i, word := range words {
		if pos, ok := index[word][doc]; ok {
			wordPositions[i] = pos
		} else {
			return false
		}
	}

	for i := range len(wordPositions) - 1 {
		for _, pos1 := range wordPositions[i] {
			for _, pos2 := range wordPositions[i+1] {
				if math.Abs(float64(pos1)-float64(pos2)) <= float64(maxDistance) {
					return true
				}
			}
		}
	}

	return false
}

func sortDocs(docs []DocRes) []DocRes {
	// cmp is reversed as I want them in descending order
	slices.SortStableFunc(docs, cmp)

	return docs
}

func cmp(a, b DocRes) int {
	if a.Occurences > b.Occurences {
		return -1
	}
	if a.Occurences < b.Occurences {
		return 1
	}
	return 0
}
