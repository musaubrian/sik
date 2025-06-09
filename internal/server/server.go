package server

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/musaubrian/sik/internal/core"
	"github.com/musaubrian/sik/internal/engine"
	"github.com/musaubrian/sik/internal/utils"
)

//go:embed www/*
var www embed.FS

type Server struct {
	port   string
	engine *engine.Engine
}

func New(index core.IndexContents) *Server {
	return &Server{
		port:   "8990",
		engine: engine.New(index),
	}
}

func (s *Server) WithPort(port string) *Server {
	s.port = port
	return s
}

func (s *Server) Start() {

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		indexHtml, err := www.ReadFile("www/index.html")
		if err != nil {
			core.Log.Error(fmt.Sprintf("[/] Failed to get index.html: %v", err))
			http.Error(w, "Unable to get page", http.StatusInternalServerError)
			return
		}
		w.Write(indexHtml)
	})
	http.HandleFunc("GET /doc", func(w http.ResponseWriter, r *http.Request) {
		markdHtml, err := www.ReadFile("www/markd.html")
		if err != nil {
			core.Log.Error(fmt.Sprintf("[/doc] Failed to get markd.html: %v", err))
			http.Error(w, "Unable to get page", http.StatusInternalServerError)
			return
		}
		w.Write(markdHtml)
	})

	http.HandleFunc("GET /js/markd.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		markMinJs, err := www.ReadFile("www/markd.min.js")
		if err != nil {
			core.Log.Error(fmt.Sprintf("[/markd.min.js] Failed to get markdjs: %v", err))
			http.Error(w, "Unable to get js asset", http.StatusInternalServerError)
			return
		}
		w.Write(markMinJs)
	})

	http.HandleFunc("POST /search", s.handleSearch)
	http.HandleFunc("GET /reload", s.handleReload)
	http.HandleFunc("GET /read-doc", s.handleReadDoc)

	core.Log.Info(fmt.Sprintf("Server running at localhost:%s", s.port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil))
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		core.Log.Error(err.Error())
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	query := r.FormValue("query")

	res, err := s.engine.Search(query)
	if err != nil {
		core.Log.Error(err.Error())
		http.Error(w, "Could not search for query", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	base, err := utils.GetSikBase()
	if err != nil {
		msg := fmt.Errorf("Could not get Base: %v", err)
		core.Log.Error(msg.Error())
		http.Error(w, msg.Error(), http.StatusInternalServerError)
		return
	}
	index, err := core.LoadIndex(utils.GetIndexLocation(base))
	if err != nil {
		msg := fmt.Errorf("Failed to load index: %v", err)
		core.Log.Error(msg.Error())
		http.Error(w, msg.Error(), http.StatusInternalServerError)
		return
	}

	s.engine = engine.New(index)
}

func (s *Server) handleReadDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fileurl := r.URL.Query().Get("url")
	f, err := os.Open(fileurl)
	if err != nil {
		msg := fmt.Errorf("Failed to open file: %v", err)
		core.Log.Error(msg.Error())
		http.Error(w, msg.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		msg := fmt.Errorf("Failed to read Files contents: %v", err)
		core.Log.Error(msg.Error())
		http.Error(w, msg.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(content.String())

	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}

}
