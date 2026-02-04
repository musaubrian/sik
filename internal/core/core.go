package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/musaubrian/logr"
	"github.com/musaubrian/sik/internal/utils"
)

var Log = logr.New()
var CurrentVersion = "v3"

type FileMeta map[string][]int // [filepath]word positions
type IndexContents map[string]FileMeta

type Index struct {
	Version  string
	Contents IndexContents
}

type MarkdownResult struct {
	Contents         map[string]string
	FilesRead        int64
	SkippedDirsCount int64
}

func ReadMarkdown(dir string) (MarkdownResult, error) {
	fileContents := map[string]string{}
	fileCount := 0
	skippedDirsCount := 0

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("COULD NOT WALK DIR: %v", err)
		}

		if d.IsDir() && utils.Ignore(d.Name()) {
			skippedDirsCount += 1
			return filepath.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			f, err := os.Open(path)
			if err != nil {
				f.Close()
				return err
			}
			defer f.Close()

			var content strings.Builder
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				content.WriteString(strings.ToLower(scanner.Text()) + " ")
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			fileContents[path] = content.String()
			fileCount += 1
		}

		return nil
	})

	return MarkdownResult{
		Contents:         fileContents,
		FilesRead:        int64(fileCount),
		SkippedDirsCount: int64(skippedDirsCount),
	}, err
}

func SaveIndex(basepath string, contents IndexContents) error {
	index := Index{Version: CurrentVersion, Contents: contents}

	json, err := json.Marshal(index)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(utils.GetIndexLocation(basepath),
		os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(json)
	return nil
}

func CreateIndex(fileContents map[string]string) (IndexContents, error) {
	index := make(IndexContents)
	for name, contents := range fileContents {
		tokenizedContents := utils.TokenizeContent(contents)
		for pos, v := range tokenizedContents {
			stemmedWord, err := utils.Stemm(v)
			if err != nil {
				return index, fmt.Errorf("Could not stem word: %v", err)
			}

			if len(stemmedWord) == 0 {
				continue
			}
			if meta, ok := index[stemmedWord]; ok {
				meta[name] = append(meta[name], pos)
			} else {
				index[stemmedWord] = FileMeta{name: []int{pos}}
			}
		}
	}

	return index, nil
}

func LoadIndex(path string) (IndexContents, error) {
	var index Index

	f, err := os.Open(path)
	if err != nil {
		return index.Contents, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&index)
	if err != nil {
		return index.Contents, fmt.Errorf("Marshalling failed: %v", err)
	}

	if index.Version != CurrentVersion {
		return index.Contents, fmt.Errorf(`Mismatched Versions: Expected <%s> Got <%s>
	Re-Index to update to the new version`,
			CurrentVersion, index.Version)
	}

	return index.Contents, nil

}
