package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/musaubrian/sik/internal/core"
	"github.com/musaubrian/sik/internal/engine"
)

//go:embed html/index.html
var html []byte

type Server struct {
	port   string
	engine *engine.Engine
}

func New(index core.Index) *Server {
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
		w.Write(html)
	})
	http.HandleFunc("POST /search", s.handleSearch)

	slog.Info(fmt.Sprintf("Server running at localhost:%s", s.port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil))
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	query := r.FormValue("query")

	res, err := s.engine.Search(query)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Could not search for query", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
