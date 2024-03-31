package strategy

import (
	"io"
	"os"
	"strings"
)

type ContentSearcher struct {
}

func (searcher *ContentSearcher) HasEntry(file string, text string) (bool, error) {
	content, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer content.Close()

	buffer := make([]byte, 64)
	for {
		n, err := content.Read(buffer)
		if err == io.EOF {
			break
		}

		if strings.Contains(string(buffer[:n]), text) {
			return true, nil
		}
	}

	return false, nil
}
