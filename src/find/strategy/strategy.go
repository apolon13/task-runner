package strategy

type Searcher interface {
	HasEntry(file string, text string) (bool, error)
}
