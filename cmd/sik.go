package mainsik

import (
	"flag"
	"fmt"
	"log/slog"
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
		contents, err := core.ReadMarkdown(indexDir)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		in, err := core.CreateIndex(contents)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		if err := core.SaveIndex(base, in); err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info("Created Index")
	}

	if browse {
		s, err := server.Bootstrap("8990")
		if err != nil {
			slog.Error(err.Error())
			return
		}
		s.Start()
	}

}
