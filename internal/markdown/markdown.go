package markdown

import (
	"bufio"
	"os"
	"strings"
)

type DocumentReader struct {
	filepath string
}

func NewDocumentReader(filepath string) *DocumentReader {
	return &DocumentReader{filepath: filepath}
}

func (r *DocumentReader) GetTitle() (string, error) {
	file, err := os.Open(r.filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case len(line) == 0:
			continue
		case len(line) == 1:
			// empty title
			return "", nil
		case line[0] == '#' && line[1] != '#':
			title := strings.TrimPrefix(line, "#")
			return title, nil
		default:
			// first none empty line is not body
			return "", nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
