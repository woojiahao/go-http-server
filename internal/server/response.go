package server

import "fmt"
import "strings"

type Response struct {
	httpVersion string
	statusCode  StatusCode
	content     string
	headers     map[string]string
}

func (r *Response) Serialize() string {
	output := make([]string, 0)
	startLine := fmt.Sprintf("%s %d %s", r.httpVersion, r.statusCode.code, r.statusCode.value)
	output = append(output, startLine)
	if len(r.headers) != 0 {
		headers := make([]string, 0)
		for key, value := range r.headers {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
		fmt.Println(headers)
		output = append(output, headers...)
	}
	output[len(output)-1] = output[len(output)-1] + "\n"

	output = append(output, r.content)
	return strings.Join(output, "\n")
}
