package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

//go:embed html/index.html
var html []byte

type Server struct {
	port string
}

func NewServer() *Server {
	return &Server{
		port: "8990",
	}
}

func (s *Server) Start() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write(html)
	})
	http.HandleFunc("POST /search", handleSearch)

	slog.Info(fmt.Sprintf("Server running at localhost:%s", s.port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	query := r.FormValue("query")

	base, err := getSikBase()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	index, err := loadIndex(getIndexLocation(base))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "load index", http.StatusInternalServerError)
		return
	}

	res, err := search(query, index)
	if err != nil {
		slog.Error(err.Error())
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
