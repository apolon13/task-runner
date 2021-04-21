package frontend

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"task-runner/cmd"
	"task-runner/config"
	"time"
)

type BuildParams struct {
	Cnf  config.Yaml
	Mode string
}

type buildProcess struct {
	dir  string
	mode string
	wait *sync.WaitGroup
	config.Command
}

func buildThreeComponentPath(a string, b string, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func Do(bp BuildParams) {
	start := time.Now()
	sem := make(chan struct{}, bp.Cnf.Build.Frontend.Parallel)
	processChannels := make(chan buildProcess)
	rootDir := bp.Cnf.Build.Frontend.Root
	fmt.Println("root directory - " + rootDir)
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	var dirNames []string
	for _, file := range files {
		if file.IsDir() {
			dirNames = append(dirNames, file.Name())
		}
	}
	go func() {
		var waitGroup sync.WaitGroup
		for _, name := range dirNames {
			filePath := buildThreeComponentPath(rootDir, name, bp.Cnf.Build.Frontend.CheckFile)
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println("add building process -" + name)
				waitGroup.Add(1)
				fp := bp.Cnf.Build.Frontend.Root + "/" + name
				cutExec := bp.Cnf.Build.Frontend.CutExecPath
				replacedFullPath := strings.ReplaceAll(fp, cutExec, "")
				processChannels <- buildProcess{
					replacedFullPath,
					bp.Mode,
					&waitGroup,
					bp.Cnf.Build.Frontend.Command,
				}
			}
			if len(bp.Cnf.Build.Frontend.Recursive) > 0 {
				for _, recursiveDir := range bp.Cnf.Build.Frontend.Recursive {
					recursive := buildThreeComponentPath(rootDir, name, recursiveDir)
					recursive = strings.ReplaceAll(recursive, recursiveDir+"/"+recursiveDir, recursiveDir)
					if _, err := os.Stat(recursive); err == nil {
						cnf := bp.Cnf
						cnf.Build.Frontend.Root = recursive
						Do(BuildParams{
							Cnf:  cnf,
							Mode: bp.Mode,
						})
					}
				}
			}
		}
		waitGroup.Wait()
		close(processChannels)
	}()

	for prc := range processChannels {
		go handlePrc(prc, sem)
	}
	close(sem)
	elapsed := time.Since(start)
	log.Printf("Elapsed time for %s - %s", bp.Cnf.Build.Frontend.Root, elapsed)
}

func handlePrc(prc buildProcess, sem chan struct{}) error {
	sem <- struct{}{}
	defer func() {
		<-sem
		prc.wait.Done()
	}()
	var args []string
	for _, arg := range prc.Command.Args {
		arg = strings.ReplaceAll(arg, "${-mode}", prc.mode)
		args = append(args, arg)
	}
	args = append(args, prc.dir)
	prc.Command.Args = args
	args = append(args, prc.dir)
	return cmd.Handle(prc.Command)
}
