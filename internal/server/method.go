package server

type Method string

const (
	GET     Method = "GET"
	PUT     Method = "PUT"
	HEAD    Method = "HEAD"
	POST    Method = "POST"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
	TRACE   Method = "TRACE"
	CONNECT Method = "CONNECT"
)

var methods = []Method{
	GET,
	PUT,
	HEAD,
	POST,
	DELETE,
	OPTIONS,
	TRACE,
	CONNECT,
}

func (m Method) isValid() bool {
	for _, method := range methods {
		if method == m {
			return true
		}
	}

	return false
}
