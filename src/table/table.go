package table

import (
	"github.com/olekukonko/tablewriter"
	"io"
)

type Table struct {
	origin *tablewriter.Table
}

func New(writer io.Writer, headers []string, borders tablewriter.Border, separator string) *Table {
	origin := tablewriter.NewWriter(writer)
	origin.SetHeader(headers)
	origin.SetBorders(borders)
	origin.SetCenterSeparator(separator)
	return &Table{origin}
}

func (table *Table) Append(item []string) {
	table.origin.Append(item)
}

func (table *Table) Render() {
	table.origin.Render()
}
