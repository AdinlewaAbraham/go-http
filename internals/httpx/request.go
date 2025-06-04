package httpx

type HttpRequest struct {
	Method    string
	URL       string
	Header    map[string]string
	Body      []byte
	Params    map[string]string
	UserAgent string
	PathParts []string
}
