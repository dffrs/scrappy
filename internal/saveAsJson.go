package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"scrappy/types"
)

const dir = "./data"

// TODO: Goal is to remove this struct
type ProdJSON struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

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

	ps := make([]ProdJSON, 0, 20)
	for _, prod := range products {
		ps = append(ps, ProdJSON{Name: prod.Name(), Price: prod.Price()})
	}

	data, err := json.MarshalIndent(ps, "", " ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
