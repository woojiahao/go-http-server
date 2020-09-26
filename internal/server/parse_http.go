package server

import (
	"fmt"
	"regexp"
	"strings"
)

func parseStartLine(startLine string) (method Method, resource string, httpVersion string, err error) {
	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {
		return Method(""), "", "", fmt.Errorf("HTTP request must include [method] [resource] [http-version]\\r\\n")
	}

	method, resource, httpVersion = Method(parts[0]), parts[1], parts[2]
	if match, _ := regexp.MatchString("^HTTP/(0.9|1.0|1.1|2.0)$", httpVersion); !match {
		err = fmt.Errorf("Invalid HTTP version. Available versions: [0.9, 1.0, 1.1, 2.0]")
		return
	}

	return
}

func parseHeader(header string) (key, value string, e error) {
	headerPattern, _ := regexp.Compile("^([\\w-_]+): *(.*)$")
	header = strings.TrimSpace(header)
	if match := headerPattern.MatchString(header); !match {
		e = fmt.Errorf("Invalid header format. Formt must be [key]: [value]")
		return
	}

	matches := headerPattern.FindStringSubmatch(header)
	key, value = matches[1], matches[2]

	return
}
