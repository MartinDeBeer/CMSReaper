package main

import (
	"encoding/json"
	"io"
	"log"
)

func DecodeJSON(data io.Reader) ([]byte, error) {
	// Create an io.Reader from the byte slice

	// Placeholder for the decoded data
	var result interface{}

	// Decode JSON data into the result variable
	if err := json.NewDecoder(data).Decode(&result); err != nil {
		log.Panic("Error decoding JSON:", err)
		return nil, err
	}

	// Optionally, re-encode the result back to []byte
	encodedResult, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return encodedResult, nil

}
