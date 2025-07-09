package main

import (
	"encoding/json"
	"os"
)

func (m model) Save() error {
	file, err := os.Create(m.fileName + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(m.events)
}
