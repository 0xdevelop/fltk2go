//go:build release

package automation

import "errors"

var ErrDisabled = errors.New("fltk2go automation debug server is disabled in release builds")

type Server struct{}

type Config struct {
	Addr string
}

func Enabled() bool { return false }

func StartDebugServer(Config) (*Server, error) { return nil, ErrDisabled }

func (s *Server) Addr() string { return "" }

func (s *Server) Close() error { return nil }
