package find

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
)

type Process struct {
	text     string
	dir      string
	searcher Searcher
	renderer Renderer
}

type processResult struct {
	filename string
	path     string
	isDir    bool
}

type Searcher interface {
	HasEntry(file string, text string) (bool, error)
}

type Renderer interface {
	AddLine(line ...string)
	Render()
}

func NewProcess(text string, dir string, searcher Searcher, renderer Renderer) Process {
	return Process{
		text,
		dir,
		searcher,
		renderer,
	}
}

func (p Process) FindEntries(threads int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errGroup, errorContext := errgroup.WithContext(ctx)
	errGroup.SetLimit(threads)
	files, err := scan(errorContext, p.dir)

	errGroup.Go(func() error {
		for {
			select {
			case <-errorContext.Done():
				return errorContext.Err()
			case file, ok := <-files:
				if !ok {
					return nil
				}
				errGroup.Go(func() error {
					hasEntry, err := p.searcher.HasEntry(file.path, p.text)
					if err == nil && hasEntry {
						p.renderer.AddLine(file.path)
					}
					return err
				})
			case scanErr := <-err:
				cancel()
				return scanErr
			}
		}
	})

	if err := errGroup.Wait(); err != nil {
		fmt.Println(err)
		return err
	}

	p.renderer.Render()
	return nil
}

func scan(ctx context.Context, root string) (chan processResult, chan error) {
	files := make(chan processResult)
	errorChan := make(chan error)
	go func() {
		defer close(files)
		defer close(errorChan)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err != nil {
					if errors.Is(err, os.ErrPermission) {
						files <- processResult{
							err.Error(),
							path,
							info.IsDir(),
						}
					} else {
						errorChan <- err
					}
				} else {
					files <- processResult{
						info.Name(),
						path,
						info.IsDir(),
					}
				}
			}
			return nil
		})
	}()

	return files, errorChan
}
