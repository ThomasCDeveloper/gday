package main

import (
	"encoding/json"
	"os"
)

func Load() ([]Event, error) {
	empty := []Event{}
	file, err := os.Open(m.fileName + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return empty, nil
		}
		return empty, err
	}
	defer file.Close()

	var loaded []Event
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loaded); err != nil {
		return empty, err
	}

	return loaded, nil
}
