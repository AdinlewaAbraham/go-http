package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net"
	"path/filepath"
	"strconv"
	"time"

	"github.com/codecrafters-io/http-server-starter-go/internals/constants"
)

type HttpResponse struct {
	conn           net.Conn
	sent           bool
	headersWritten bool
	headers        map[string]string
	status         int
	statusText     string
	bodyBuffer     *bytes.Buffer
	protocol       string
}

func (res *HttpResponse) writeChunk(data []byte) {
	chunkSize := fmt.Sprintf("%x\r\n", len(data))
	res.conn.Write([]byte(chunkSize))

	res.conn.Write(data)
	res.conn.Write([]byte("\r\n"))
}

func (res *HttpResponse) fillDefaults() {
	if res.protocol == "" {
		res.protocol = "HTTP/1.1"
	}

	if res.headers == nil {
		res.headers = make(map[string]string)
	}

	if res.status == 0 {
		res.status = 200
	}

	if res.statusText == "" {
		res.statusText = constants.StatusTexts[constants.StatusCode(res.status)]
	}

	if _, exists := res.headers["Content-Length"]; !exists &&
		res.bodyBuffer.Len() > 0 &&
		res.headers["Transfer-Encoding"] != "chunked" {
		res.headers["Content-Length"] = strconv.Itoa(res.bodyBuffer.Len())
	}

	if _, exists := res.headers["Date"]; !exists {
		res.headers["Date"] = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	}

	if _, exists := res.headers["Server"]; !exists {
		res.headers["Server"] = "CustomServer/1.0"
	}
}

func (res *HttpResponse) buildHTTPResponse() []byte {
	var buf bytes.Buffer

	res.fillDefaults()

	estimatedSize := 512 + res.bodyBuffer.Len()
	buf.Grow(estimatedSize)

	buf.WriteString(res.protocol)
	buf.WriteByte(' ')
	buf.WriteString(strconv.Itoa(int(res.status)))
	buf.WriteByte(' ')
	buf.WriteString(res.statusText)
	buf.WriteString("\r\n")

	for name, value := range res.headers {
		buf.WriteString(name)
		buf.WriteString(": ")
		buf.WriteString(value)
		buf.WriteString("\r\n")
	}

	buf.WriteString("\r\n")

	if res.bodyBuffer.Len() > 0 {
		buf.Write(res.bodyBuffer.Bytes())
	}

	return buf.Bytes()
}

func getStatusText(code constants.StatusCode) string {
	if text, exists := constants.StatusTexts[code]; exists {
		return text
	}
	return "Unknown"
}

func NewResponse(conn net.Conn) (res *HttpResponse) {
	return &HttpResponse{
		conn:       conn,
		status:     200,
		headers:    make(map[string]string),
		statusText: "",
		bodyBuffer: &bytes.Buffer{},
	}

}

func (res *HttpResponse) SetHeader(key, value string) {
	if res.headersWritten {
		fmt.Println("header already written cannot write key: %s and value: %s", key, value)
		return
	}
	res.headers[key] = value
}

func (res *HttpResponse) WriteHeader(status int, header map[string]string) {
	if res.headersWritten {
		fmt.Println("header already written cannot write status: %s", status)
		return
	}
	if res.protocol == "" {
		res.protocol = "HTTP/1.1"
	}
	if header != nil {
		for k, v := range header {
			res.headers[k] = v
		}
	}
	res.status = status
	res.headersWritten = true
}

func (res *HttpResponse) Write(data []byte) {
	if res.sent {
		fmt.Println("Already flushed res")
		return
	}

	if res.headers["Transfer-Encoding"] == "chunked" {
		res.writeChunk(data)
	} else {
		res.bodyBuffer.Write(data)
	}
}

func (res *HttpResponse) Send(body []byte) error {
	if res.sent {
		return errors.New("already sent a response")
	}

	if _, exists := res.headers["Content-Type"]; !exists && len(body) > 0 {
		res.headers["Content-Type"] = "text/plain"
	}

	res.bodyBuffer.Write(body)
	res.Flush()
	return nil
}

func (res *HttpResponse) SendFile(filename string, body []byte) error {
	if res.sent {
		return errors.New("response already sent")
	}

	if _, exists := res.headers["Content-Type"]; !exists {
		ext := filepath.Ext(filename)
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			res.headers["Content-Type"] = contentType
		} else {
			res.headers["Content-Type"] = "application/octet-stream"
		}
	}

	res.bodyBuffer.Write(body)
	res.Flush()
	return nil
}

func (res *HttpResponse) SendHTML(body []byte) error {
	res.headers["Content-Type"] = "text/html; charset=utf-8"
	return res.Send(body)
}

func (res *HttpResponse) SendJSON(data map[string]any) error {
	s, err := json.Marshal(data)

	if err != nil {
		return errors.New("Could not encode json")
	}

	res.headers["Content-Type"] = "application/json"
	res.bodyBuffer.Write(s)
	res.Flush()

	return nil
}

func (res *HttpResponse) Status(code int) *HttpResponse {
	res.status = code
	res.statusText = getStatusText(constants.StatusCode(code))
	return res
}

func (res *HttpResponse) StatusText(text string) *HttpResponse {
	res.statusText = text
	return res
}

func (res *HttpResponse) End() error {
	err := res.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (res *HttpResponse) Flush() error {
	response := res.buildHTTPResponse()
	_, err := res.conn.Write(response)
	if err != nil {
		return err
	}
	res.sent = true
	err = res.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
