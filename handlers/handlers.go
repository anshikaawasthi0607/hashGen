//this file sets up routes, handles req
package handlers

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"hashservice/hashgen"
)

type Handler struct {
	tmpl *template.Template
}

type HashRequest struct {
	Input string `json:"input"`
}

type HashResponse struct {
	Input string `json:"input"`
	Hash  string `json:"hash"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SetupRoutes(templatesFS fs.FS, staticFS fs.FS) http.Handler {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html") 
	if err != nil {
		log.Fatalf("Failed to parse embedded templates: %v", err)
	}

	h := &Handler{tmpl: tmpl}

	mux := http.NewServeMux() 

	mux.HandleFunc("/", h.Index)
	
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Hash generation API
	mux.HandleFunc("/api/hash", h.GenerateHash) //frontend calls

	// Liveness probe
	mux.HandleFunc("/health", h.HealthCheck)

	return requestLogger(mux)
}


func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) GenerateHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{"Method not allowed"})
		return
	}

	var req HashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{"Invalid JSON body"})
		return
	}

	hashValue, err := hashgen.Generate(req.Input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, HashResponse{
		Input: req.Input,
		Hash:  hashValue,
	})
}


func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

// a middleware that logs each incoming request.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}