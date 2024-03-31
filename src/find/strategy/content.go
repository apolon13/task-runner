package strategy

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type ContentSearcher struct {
}

func (searcher *ContentSearcher) HasEntry(file string, text string) (bool, error) {
	content, err := os.Open(file)
	defer content.Close()
	if err != nil {
		return false, err
	}

	fileInfo, err := content.Stat()
	if err != nil {
		return false, err
	}

	if fileInfo.IsDir() {
		return false, nil
	}

	scanner := bufio.NewReader(content)
	for {
		res, err := scanner.ReadString('\n')
		if err == io.EOF {
			break
		}
		if strings.Contains(res, text) {
			return true, nil
		}
	}

	return false, nil
}
