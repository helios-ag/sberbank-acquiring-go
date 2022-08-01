package testing

import (
	"net/http"
	"net/http/httptest"
)

// Server Test Environment
type Server struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	URL    string
}

func (server *Server) Teardown() {
	server.Server.Close()
	server.Server = nil
	server.Mux = nil
}

func NewServer() Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return Server{
		Server: server,
		Mux:    mux,
		URL:    server.URL,
	}
}
