package find

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
	"os"
	"sync"
	"task-runner/src/find/strategy"
	"task-runner/src/table"
)

type findProcess struct {
	text           string
	dir            string
	findInFilename bool
	subProcesses   chan *findProcess
	wg             *sync.WaitGroup
	isRoot         bool
	strategy       strategy.Searcher
}

type processResult struct {
	filename string
	error    error
	dir      string
	isDir    bool
}

var Find = &cobra.Command{
	Use:   "find",
	Short: "find text in filename or content",
	RunE: func(cmd *cobra.Command, args []string) error {
		var wg sync.WaitGroup
		path, _ := cmd.Flags().GetString("path")
		text, _ := cmd.Flags().GetString("text")
		strategyType, _ := cmd.Flags().GetString("strategy")
		threads, _ := cmd.Flags().GetInt("threads")

		var searchStrategy strategy.Searcher
		if strategyType == "filename" {
			searchStrategy = &strategy.FilenameSearcher{}
		} else if strategyType == "content" {
			searchStrategy = &strategy.ContentSearcher{}
		}

		process := &findProcess{
			text,
			path,
			true,
			make(chan *findProcess),
			&wg,
			true,
			searchStrategy,
		}
		return process.findEntries(int64(threads))
	},
}

func init() {
	Find.Flags().StringP("path", "p", "/", "root dir")
	Find.Flags().StringP("text", "t", "", "search text")
	Find.Flags().StringP("strategy", "s", "filename", "search strategy (filename|content)")
	Find.Flags().IntP("threads", "c", 5, "count threads")
	Find.MarkFlagRequired("text")
}

func (process *findProcess) findEntries(threads int64) error {
	sem := semaphore.NewWeighted(threads)
	output := make(chan processResult)
	resultTable := table.New(
		os.Stdout,
		[]string{"type", "name", "dir"},
		tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false},
		"|",
	)

	go func() {
		process.scan(sem, output)
		process.wg.Wait()
		resultTable.Render()
		close(process.subProcesses)
	}()

	go func() {
		for {
			select {
			case result := <-output:
				if result.error != nil {
					_ = fmt.Errorf("%s", result.error)
				} else {
					var itemType string
					if result.isDir {
						itemType = "DIR"
					} else {
						itemType = "FILE"
					}
					resultTable.Append([]string{
						itemType,
						result.filename,
						result.dir,
					})
				}
			}
		}
	}()

	for subProcess := range process.subProcesses {
		go subProcess.scan(sem, output)
	}
	return nil
}

func (process *findProcess) scan(semaphore *semaphore.Weighted, output chan processResult) {
	if !process.isRoot {
		semaphore.Acquire(context.TODO(), 1)
	}

	res, err := os.ReadDir(process.dir)
	defer func() {
		if !process.isRoot {
			semaphore.Release(1)
			process.wg.Done()
		}
	}()

	if err != nil {
		output <- processResult{"", err, "", false}
	}

	for _, item := range res {
		fullPath := process.dir + "/" + item.Name()
		if item.IsDir() {
			process.wg.Add(1)
			process.subProcesses <- &findProcess{
				process.text,
				fullPath,
				process.findInFilename,
				process.subProcesses,
				process.wg,
				false,
				process.strategy,
			}
		}

		hasEntry, err := process.strategy.HasEntry(fullPath, process.text)
		if err != nil {
			output <- processResult{"", err, "", false}
		}

		if hasEntry {
			output <- processResult{
				item.Name(),
				nil,
				process.dir,
				item.IsDir(),
			}
		}
	}
}
