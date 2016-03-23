package web

import "net/http"

// Server implements a web server for playing Go
type Server struct {
	Public string
}

// Bind configures routes in the provided http.ServeMux
func (s *Server) Bind(mux *http.ServeMux) error {
	mux.Handle("/", http.FileServer(http.Dir(s.Public)))
	return nil
}
