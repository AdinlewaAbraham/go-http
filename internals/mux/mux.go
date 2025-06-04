package mux

import (
	"fmt"
	"log"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internals/httpx"
)

type HttpMuxTrieNode struct {
	segment     string
	handler     map[string]func(request *httpx.HttpRequest, response *httpx.HttpResponse)
	middlewares []func(req *httpx.HttpRequest, res *httpx.HttpResponse, next func())

	staticChildren map[string]*HttpMuxTrieNode

	paramChild *HttpMuxTrieNode
	paramName  string
}

type HttpMux struct {
	muxTrieRoot *HttpMuxTrieNode
}

type RouteMatch struct {
	Handler     func(*httpx.HttpRequest, *httpx.HttpResponse)
	Middlewares []func(*httpx.HttpRequest, *httpx.HttpResponse, func())
	Params      map[string]string
	Found       bool
}

func chainMiddleware(handler func(request *httpx.HttpRequest, response *httpx.HttpResponse), middlewares []func(req *httpx.HttpRequest, res *httpx.HttpResponse, next func())) func(request *httpx.HttpRequest, response *httpx.HttpResponse) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		current := middlewares[i]
		next := handler
		handler = func(req *httpx.HttpRequest, res *httpx.HttpResponse) {
			current(req, res, func() {
				next(req, res)
			})
		}
	}
	return handler
}

func extractParam(segment string) (paramName string, isParam bool) {
	isParam = strings.HasPrefix(segment, ":")
	if isParam && len(segment) > 1 {
		paramName = segment[1:]
	} else if isParam {
		isParam = false
	}
	return paramName, isParam
}

func newMuxTrieNode(pathSegment string) *HttpMuxTrieNode {
	return &HttpMuxTrieNode{
		segment:        pathSegment,
		staticChildren: make(map[string]*HttpMuxTrieNode),
		handler:        make(map[string]func(request *httpx.HttpRequest, response *httpx.HttpResponse)),
	}
}

func NewHttpMux() *HttpMux {
	return &HttpMux{
		muxTrieRoot: newMuxTrieNode("*"),
	}
}

func (m *HttpMux) GetMuxTrieRoot() *HttpMuxTrieNode {
	return m.muxTrieRoot
}

func (m *HttpMux) ExplorePath(method string, segments []string) RouteMatch {
	node := m.muxTrieRoot
	match := RouteMatch{
		Middlewares: append([]func(*httpx.HttpRequest, *httpx.HttpResponse, func()){}, node.middlewares...),
		Params:      map[string]string{},
	}

	for _, segment := range segments {
		if next, ok := node.staticChildren[segment]; ok {
			node = next
			match.Middlewares = append(match.Middlewares, node.middlewares...)
		} else if node.paramChild != nil {
			log.Printf("node", node)
			node = node.paramChild
			log.Printf("node paramchild", node)
			log.Printf("writing key: %s to value: ", node.paramName, segment)
			match.Params[node.paramName] = segment
			match.Middlewares = append(match.Middlewares, node.middlewares...)
		} else {
			return match
		}
	}

	handler := node.handler[method]
	if handler == nil {
		return match
	}

	match.Handler = handler
	match.Found = true
	return match
}

func (m *HttpMux) RegisterRoute(path string, method string, handler func(request *httpx.HttpRequest, response *httpx.HttpResponse)) {
	// validate path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	currentNode := m.muxTrieRoot

	for _, pathSegment := range pathParts {
		paramName, isParam := extractParam(pathSegment)
		if isParam {
			log.Printf("registering route with path name %s gotten from segment %s", paramName, pathSegment)
			child := currentNode.paramChild
			if child == nil {
				child = newMuxTrieNode(pathSegment)
				child.paramName = paramName
				currentNode.paramChild = child
			}
			currentNode = child
		} else {
			if currentNode.staticChildren == nil {
				currentNode.staticChildren = make(map[string]*HttpMuxTrieNode)
			}

			child, exists := currentNode.staticChildren[pathSegment]
			if !exists {
				child = newMuxTrieNode(pathSegment)
				currentNode.staticChildren[pathSegment] = child
			}
			currentNode = child
		}
	}
	currentNode.handler[method] = handler
}

func (m *HttpMux) RouteRequest(req *httpx.HttpRequest, res *httpx.HttpResponse) {
	match := m.ExplorePath(req.Method, req.PathParts)

	if !match.Found {
		fmt.Println("Handler not found for", req.Method, req.URL)
		res.Status(404).Send([]byte("page not found"))
		return
	}

	log.Printf("matched params: ", match.Params)
	req.Params = match.Params

	finalHandler := chainMiddleware(match.Handler, match.Middlewares)
	finalHandler(req, res)
}

func (m *HttpMux) AttachMiddleware(path string, mw func(req *httpx.HttpRequest, res *httpx.HttpResponse, next func())) {
	pathParts := strings.Split(path, "/")
	currentNode := m.muxTrieRoot
	for _, pathSegment := range pathParts {
		_, isParam := extractParam(pathSegment)
		if isParam {
			if currentNode.paramChild == nil {
				fmt.Println("(Node not found) Could not attach middleware to path: ", path)
				return
			}
			currentNode = currentNode.paramChild
		} else {
			val, exists := currentNode.staticChildren[pathSegment]
			if !exists {
				fmt.Println("(Node not found) Could not attach middleware to path: ", path)
				return
			}
			currentNode = val
		}
	}
	currentNode.middlewares = append(currentNode.middlewares, mw)
}
