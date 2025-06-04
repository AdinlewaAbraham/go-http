package util

import (
	"bufio"
	"errors"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internals/constants"
)

type ParsedRequestInfo struct {
	Method    string
	Path      string
	Host      string
	UserAgent string
}

func ExtractUrl(buff []byte) (*ParsedRequestInfo, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(buff)))

	var info ParsedRequestInfo

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

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

		if strings.HasPrefix(line, "Host: ") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				info.Host = strings.TrimSpace(parts[1])
			}
		}

		if strings.HasPrefix(line, "User-Agent: ") {
			info.UserAgent = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent: "))
		}
	}

	if info.Path == "" {
		return nil, errors.New("could not find path line")
	}

	if info.Host == "" {
		return nil, errors.New("could not find host line")
	}

	return &info, nil
}
