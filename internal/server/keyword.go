package server

type Keyword string

const (
	GET    Keyword = "GET"
	ANSWER Keyword = "ANSWER"
	ERROR  Keyword = "ERROR"
	SET    Keyword = "SET"
	CLEAR  Keyword = "CLEAR"
	ALL    Keyword = "ALL"
)
