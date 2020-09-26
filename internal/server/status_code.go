package server

type StatusCode struct {
	code  int
	value string
}

var (
	OK                  StatusCode = StatusCode{200, "OK"}
	BadRequest          StatusCode = StatusCode{400, "Bad Request"}
	NotFound            StatusCode = StatusCode{404, "Not Found"}
	InternalServerError StatusCode = StatusCode{500, "Internal Server Error"}
)
