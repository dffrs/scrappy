package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"scrappy/types"
)

const dir = "./data"

func createDirIfNotExists() (string, error) {
	path := filepath.Join(".", dir)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path, nil
}

func SaveAsJSON(products []types.Product, fileName string) error {
	dir, err := createDirIfNotExists()
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s%s%s%s", dir, string(os.PathSeparator), fileName, ".json"))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	data, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
