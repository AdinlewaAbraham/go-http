package httpx

import "github.com/codecrafters-io/http-server-starter-go/internals/util"

type HttpRequest struct {
	util.ParsedRequestInfo
	URL       string
	Params    map[string]string
	PathParts []string
}
