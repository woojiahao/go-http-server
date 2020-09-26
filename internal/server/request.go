package server

type Request struct {
	method      Method
	resource    string
	httpVersion string
	headers     map[string]string
}
