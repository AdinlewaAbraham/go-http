package server

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internals/httpx"
	"github.com/codecrafters-io/http-server-starter-go/internals/mux"
	"github.com/codecrafters-io/http-server-starter-go/internals/util"
)

type RouteMethod func(path string, handler *func(request *httpx.HttpRequest, response *httpx.HttpResponse))
type Server struct {
	mux mux.HttpMux
}

func (s *Server) Listen(port string, cb func()) error {
	l, err := net.Listen("tcp", port)

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to bind to port %s %s", port, err.Error())
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	defer l.Close()

	fmt.Println("Server listening on port 4221...")
	go cb()

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go func() {
			defer conn.Close()

			req, err := s.parseConn(conn)
			if err != nil {
				// shit went wrong
				return
			}
			res := httpx.NewResponse(conn)

			s.mux.RouteRequest(req, res)
		}()
	}

}

func (s *Server) Use(path string, mw func(req *httpx.HttpRequest, res *httpx.HttpResponse, next func())) {
	s.mux.AttachMiddleware(path, mw)
}

func (s *Server) parseConn(conn net.Conn) (req *httpx.HttpRequest, err error) {
	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		return req, err
	}
	
	info, extractUrlErr := util.ExtractUrl(buff)
	if extractUrlErr != nil {
		return req, extractUrlErr
	}

	segments := strings.Split(strings.Trim(info.Path, "/"), "/")

	req = &httpx.HttpRequest{
		Method:    info.Method,
		URL:       info.Host + info.Path,
		UserAgent: info.UserAgent,
		PathParts: segments,
	}
	return req, nil
}

func CreateServer() *Server {
	return &Server{
		mux: *mux.NewHttpMux(),
	}
}
