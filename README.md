# go-http-server
HTTP server built in Go, following this guide: https://theprogrammershangout.com/resources/projects/http-project-guide/intro.md

## Run server

To run the server, use `go run cmd/main.go`.

## Connect to the server

You can connect to the server using the terminal via `netcat -v -v localhost 8000`.

## Protocols

- `GET` - gets a word from the dictionary. Returns error if not found
- `SET` - add a new word and its definition to the dictionary
- `CLEAR` - clears the dictionary
- `ALL` - returns all of the available words in the dictionary
