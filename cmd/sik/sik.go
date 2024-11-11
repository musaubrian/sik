package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/musaubrian/sik/internal/core"
	"github.com/musaubrian/sik/internal/server"
	"github.com/musaubrian/sik/internal/utils"
)

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

	base, err := utils.GetSikBase()
	if err != nil {
		core.Logging.Error(err.Error())
		return
	}

	if _, err := os.Stat(base); err != nil {
		if err := os.Mkdir(base, 0755); err != nil {
			core.Logging.Error(err.Error())
			return
		}
		core.Logging.Info(fmt.Sprintf("Created %s", filepath.Base(base)))
	}

	if len(indexDir) > 0 {
		contents, err := core.ReadMarkdown(indexDir)
		if err != nil {
			core.Logging.Error(err.Error())
			return
		}

		in, err := core.CreateIndex(contents)
		if err != nil {
			core.Logging.Error(err.Error())
			return
		}
		if err := core.SaveIndex(base, in); err != nil {
			core.Logging.Error(err.Error())
			return
		}
		core.Logging.Info("Created Index")
	}

	if browse {
		s, err := serverBootstrap()
		if err != nil {
			core.Logging.Error(err.Error())
			return
		}
		s.Start()
	}

}

func serverBootstrap() (*server.Server, error) {
	base, err := utils.GetSikBase()
	if err != nil {
		return nil, fmt.Errorf("failed to get sik base: %w", err)
	}

	index, err := core.LoadIndex(utils.GetIndexLocation(base))
	if err != nil {
		return nil, fmt.Errorf("failed to load index: %w", err)
	}

	return server.New(index).WithPort("4000"), nil
}
