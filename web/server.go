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

	return nil
}

// Bind configures routes in the provided http.ServeMux
func (s *Server) Bind(mux *http.ServeMux) error {
	mux.Handle("/board.json", s.handler(s.serveBoard))
	mux.Handle("/move", s.handler(s.handleMove))
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

func replyJSON(w http.ResponseWriter, val interface{}) error {
	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Printf("error encoding json: %v", err)
		return err
	}
	return nil
}

func (s *Server) handler(handler func(http.ResponseWriter, *http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, err := handler(w, r); err != nil {
			log.Printf("error: %v", err)
			if ue, ok := err.(*UserError); ok {
				w.WriteHeader(ue.Code())
				replyJSON(w, ue)
			} else {
				w.WriteHeader(500)
				replyJSON(w, &struct {
					Err string `json:"error"`
				}{"an internal error occurred"})
			}
		} else {
			replyJSON(w, val)
		}
	})
}

func (s *Server) serveBoard(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	return &out, nil
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var args struct {
		X      int    `json:"x"`
		Y      int    `json:"y"`
		ToMove string `json:"to_move"`
	}
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return nil, &UserError{err.Error()}
	}

	if args.ToMove != colorStr(s.game.ToPlay()) {
		return nil, &UserError{"it's not your turn"}
	}

	if err := s.game.Move(args.X, args.Y); err != nil {
		return nil, &UserError{"illegal move"}
	}

	return s.serveBoard(w, r)
}
