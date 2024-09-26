package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FileMeta map[string][]int // [filepath]word positions

type Index map[string]FileMeta

func main() {
	var (
		indexDir string
		browse   bool
	)

	flag.StringVar(&indexDir, "index", "", "Path to directory to index")

	flag.BoolVar(&browse, "b", false, "Start up the searching page")
	flag.Parse()

	if len(indexDir) == 0 && !browse {
		flag.Usage()
		return
	}

	base, err := getSikBase()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if _, err := os.Stat(base); err != nil {
		if err := os.Mkdir(base, 0755); err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info(fmt.Sprintf("Created %s", filepath.Base(base)))
	}

	if len(indexDir) > 0 {
		contents, err := readMarkdown(indexDir)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		in, err := createIndex(contents)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if err := saveIndex(base, in); err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info("Created Index")
	}

	s := NewServer()

	if browse {
		s.Start()
	}

}

func saveIndex(basepath string, index Index) error {
	json, err := json.Marshal(index)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(getIndexLocation(basepath), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(json)
	return nil
}

func createIndex(fileContents map[string]string) (Index, error) {
	index := make(Index)
	for name, contents := range fileContents {
		tokenizedContents := tokenizeContent(contents)
		for pos, v := range tokenizedContents {
			stemmedWord, err := stemm(v)
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

func loadIndex(path string) (Index, error) {
	index := make(Index)

	f, err := os.Open(path)
	if err != nil {
		return index, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&index)
	if err != nil {
		return index, fmt.Errorf("Marshalling failed: %v", err)
	}

	return index, nil

}

func readMarkdown(dir string) (map[string]string, error) {
	fileContents := map[string]string{}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("COULD NOT WALK DIR: %v", err)
		}

		if d.IsDir() && ignore(d.Name()) {
			slog.Info("SKIPPING " + d.Name())
			return filepath.SkipDir
		}

		// TODO: add more extensions maybe, will think about this
		// Will need to change how checks happen
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") || strings.HasSuffix(d.Name(), ".txt") {
			f, err := os.Open(path)
			if err != nil {
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
		}

		return nil
	})

	return fileContents, err
}
