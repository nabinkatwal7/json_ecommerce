package database

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var mu sync.Mutex

func EnsureFile(filePath string) error {
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(filePath, []byte("{}"), 0644)
	}

	return nil
}

func ReadJSON[T any](filePath string)([]T, error) {
	mu.Lock()
	defer mu.Unlock()

	if err := EnsureFile(filePath); err != nil {
		return nil, err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []T

	if len(file) == 0 {
		return []T{}, nil
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func WriteJSON[T any](filePath string, data []T) error {
	mu.Lock()
	defer mu.Unlock()

	if err := EnsureFile(filePath); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}