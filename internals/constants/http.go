package constants

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
)

var AllHTTPMethods = []string{
	GET,
	POST,
	PUT,
	DELETE,
	PATCH,
	OPTIONS,
	HEAD,
}

type StatusCode int

const (
	StatusOK        StatusCode = 200
	StatusCreated   StatusCode = 201
	StatusAccepted  StatusCode = 202
	StatusNoContent StatusCode = 204

	StatusMovedPermanently StatusCode = 301
	StatusFound            StatusCode = 302
	StatusNotModified      StatusCode = 304

	StatusBadRequest          StatusCode = 400
	StatusUnauthorized        StatusCode = 401
	StatusForbidden           StatusCode = 403
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusConflict            StatusCode = 409
	StatusUnprocessableEntity StatusCode = 422

	StatusInternalServerError StatusCode = 500
	StatusNotImplemented      StatusCode = 501
	StatusBadGateway          StatusCode = 502
	StatusServiceUnavailable  StatusCode = 503
)

var StatusTexts = map[StatusCode]string{
	StatusOK:        "OK",
	StatusCreated:   "Created",
	StatusAccepted:  "Accepted",
	StatusNoContent: "No Content",

	StatusMovedPermanently: "Moved Permanently",
	StatusFound:            "Found",
	StatusNotModified:      "Not Modified",

	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusConflict:            "Conflict",
	StatusUnprocessableEntity: "Unprocessable Entity",

	StatusInternalServerError: "Internal Server Error",
	StatusNotImplemented:      "Not Implemented",
	StatusBadGateway:          "Bad Gateway",
	StatusServiceUnavailable:  "Service Unavailable",
}

type ContentType string

const (
	ContentTypeJSON  ContentType = "application/json"
	ContentTypeHTML  ContentType = "text/html"
	ContentTypeText  ContentType = "text/plain"
	ContentTypeOctet ContentType = "application/octet-stream"
)

const (
	HeaderKeyContentType   ContentType = "Content-Type"
	HeaderKeyContentLength ContentType = "Content-Length"
)
