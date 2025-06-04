package util

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internals/constants"
)

type ParsedRequestInfo struct {
	Method        string
	Path          string
	Host          string
	UserAgent     string
	ContentType   string
	ContentLength int
	Header        map[string]string
	Body          []byte
}

func ExtractUrl(buff []byte) (*ParsedRequestInfo, error) {
	parts := bytes.SplitN(buff, []byte("\r\n\r\n"), 2)
	if len(parts) < 2 {
		return nil, errors.New("invalid HTTP request format")
	}

	headersPart := string(parts[0])
	bodyPart := parts[1]

	info := ParsedRequestInfo{
		Header: make(map[string]string),
	}
	info.Body = bodyPart

	scanner := bufio.NewScanner(strings.NewReader(headersPart))
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		if info.Method == "" {
			for _, m := range constants.AllHTTPMethods {
				if strings.HasPrefix(line, m+" ") {
					parts := strings.Split(line, " ")
					if len(parts) >= 2 {
						info.Method = m
						info.Path = strings.TrimSpace(parts[1])
					}
					break
				}
			}
		}

		lineParts := strings.SplitN(line, ":", 2)
		if len(lineParts) < 2 {
			continue
		}

		header := strings.TrimSpace(lineParts[0])
		content := strings.TrimSpace(lineParts[1])

		switch header {
		case "Host":
			info.Host = content
		case "User-Agent":
			info.UserAgent = content
		case "Content-Type":
			info.ContentType = content
		case "Content-Length":
			length, err := strconv.Atoi(content)
			if err != nil {
				log.Printf("invalid Content-Length: %v", err)
			} else {
				info.ContentLength = length
			}
		}

		info.Header[header] = content
	}

	if info.Path == "" {
		return nil, errors.New("could not find path line")
	}

	if info.Host == "" {
		return nil, errors.New("could not find host line")
	}

	return &info, nil
}
