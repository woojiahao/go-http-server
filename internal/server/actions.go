package server

import "fmt"
import "strings"

var data = map[string]string{
	"hello": "world",
	"open":  "sesame",
}

func formatData() string {
	output := make([]string, 0)
	for key, value := range data {
		output = append(output, fmt.Sprintf("%s: %s", key, value))
	}
	return strings.Join(output, "\n")
}

func handleGET(word string) (string, error) {
	for key, value := range data {
		if key == word {
			return value, nil
		}
	}
	return "", fmt.Errorf("can't find %s", word)
}

func handleSET(word, value string) {
	data[word] = value
}

func handleCLEAR() {
	data = make(map[string]string)
}

func handleALL() []string {
	words := make([]string, 0)
	for key, _ := range data {
		words = append(words, key)
	}
	return words
}
