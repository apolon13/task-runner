package main

import (
	"github.com/spf13/cobra"
	"task-runner/src/cmd/find"
)

var root = &cobra.Command{
	Use:   "task-runner",
	Short: "task-runner is multi command util",
}

func init() {
	root.AddCommand(find.Find)
}

func main() {
	root.Execute()
}
