package openapi

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type decoder interface {
	Decode(any) error
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Document struct {
	Info Info `json:"info"`
}

type DocumentReader struct {
	filepath string
}

func NewDocumentReader(filepath string) *DocumentReader {
	return &DocumentReader{filepath: filepath}
}

func (r *DocumentReader) GetTitle() (string, error) {
	f, err := os.Open(r.filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var dec decoder
	switch filepath.Ext(r.filepath) {
	case ".json":
		dec = json.NewDecoder(f)
	case ".yaml", ".yml":
		dec = yaml.NewDecoder(f)
	default:
		return "", errors.New("unsupported file type")
	}

	doc := Document{}
	if err := dec.Decode(&doc); err != nil {
		return "", err
	}

	return doc.Info.Title, nil
}
