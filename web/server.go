package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"nelhage.com/minigo/game"
)

// DefaultSize is the default board size if none is provided
const DefaultSize = 9

// Config represents the configuration for a server
type Config struct {
	Public string
	Size   int
}

// Server implements a web server for playing Go
type Server struct {
	c Config

	game *game.Game
}

// Init configures a server and initializes any relevant
func (s *Server) Init(c *Config) error {
	s.c = *c
	size := s.c.Size
	if size == 0 {
		size = DefaultSize
	}
	s.game = game.New(size)

	s.game.Move(4, 4)
	s.game.Move(3, 3)
	s.game.Move(6, 6)
	return nil
}

// Bind configures routes in the provided http.ServeMux
func (s *Server) Bind(mux *http.ServeMux) error {
	mux.Handle("/board.json", http.HandlerFunc(s.serveBoard))
	mux.Handle("/", http.FileServer(http.Dir(s.c.Public)))
	return nil
}

func colorStr(c game.Color) string {
	switch c {
	case game.White:
		return "W"
	case game.Black:
		return "B"
	default:
		panic(fmt.Sprintf("bad color %v", c))
	}
}

func (s *Server) serveBoard(w http.ResponseWriter, r *http.Request) {
	var out struct {
		ToMove    string            `json:"to_move"`
		Positions map[string]string `json:"positions"`
	}
	out.Positions = make(map[string]string)
	for x := 0; x < s.game.Size; x++ {
		for y := 0; y < s.game.Size; y++ {
			if c, ok := s.game.At(x, y); ok {
				key := fmt.Sprintf("%d,%d", x, y)
				out.Positions[key] = colorStr(c)
			}
		}
	}

	out.ToMove = colorStr(s.game.ToPlay())

	if err := json.NewEncoder(w).Encode(&out); err != nil {
		log.Printf("error encoding json: %v", err)
	}
}
