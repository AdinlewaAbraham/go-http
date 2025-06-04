package server

import (
	"github.com/codecrafters-io/http-server-starter-go/internals/httpx"
)

func (s *Server) Get(path string, handler func(request *httpx.HttpRequest, response *httpx.HttpResponse)) {
	s.mux.RegisterRoute(path, "GET", handler)

}
func (s *Server) Post(path string, handler func(request *httpx.HttpRequest, response *httpx.HttpResponse)) {
	s.mux.RegisterRoute(path, "POST", handler)
}
func (s *Server) Put(path string, handler func(request *httpx.HttpRequest, response *httpx.HttpResponse)) {
	s.mux.RegisterRoute(path, "PUT", handler)
}
func (s *Server) Delete(path string, handler func(request *httpx.HttpRequest, response *httpx.HttpResponse)) {
	s.mux.RegisterRoute(path, "DELETE", handler)
}
