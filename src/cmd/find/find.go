package find

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"task-runner/src/proc/find"
	"task-runner/src/search/strategy"
	"task-runner/src/view/table"
)

var Find = &cobra.Command{
	Use:   "find",
	Short: "find text in filename or content",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		text, _ := cmd.Flags().GetString("text")
		threads, _ := cmd.Flags().GetInt("threads")
		searcherType, _ := cmd.Flags().GetString("strategy")

		var searcher find.Searcher
		if searcherType == "filename" {
			searcher = strategy.FilenameSearcher{}
		} else if searcherType == "content" {
			searcher = strategy.ContentSearcher{}
		}

		proc := find.NewProcess(text, path, searcher, table.New(
			os.Stdout,
			[]string{"Path"},
			tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false},
			"|",
		))
		return proc.FindEntries(threads)
	},
}

func init() {
	Find.Flags().StringP("path", "p", "/", "root dir")
	Find.Flags().StringP("text", "t", "", "search text")
	Find.Flags().StringP("strategy", "s", "filename", "search strategy (filename|content)")
	Find.Flags().IntP("threads", "c", 50, "count threads")
	Find.MarkFlagRequired("text")
}
