package table

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"sync"
)

type Table struct {
	origin *tablewriter.Table
	mu     sync.Mutex
}

func New(writer io.Writer, headers []string, borders tablewriter.Border, separator string) *Table {
	origin := tablewriter.NewWriter(writer)
	origin.SetHeader(headers)
	origin.SetBorders(borders)
	origin.SetCenterSeparator(separator)
	return &Table{origin, sync.Mutex{}}
}

func (table *Table) AddLine(line ...string) {
	table.mu.Lock()
	defer table.mu.Unlock()
	table.origin.Append(line)
}

func (table *Table) Render() {
	table.origin.Render()
}
