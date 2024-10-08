package strategy

import "strings"

type FilenameSearcher struct {
}

func (searcher FilenameSearcher) HasEntry(file string, text string) (bool, error) {
	return strings.Contains(file, text), nil
}
