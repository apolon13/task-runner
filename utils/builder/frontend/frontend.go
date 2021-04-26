package frontend

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"task-runner/cmd"
	"time"
)

type BuildProcess struct {
	Mode          string
	ProcessParams *ProcessParams
}

type ProcessParams struct {
	Root        string   `yaml:"root"`
	CutExecPath string   `yaml:"cut-exec-path"`
	Parallel    int      `yaml:"parallel"`
	Recursive   []string `yaml:"recursive"`
	CheckFile   string   `yaml:"check-file"`
	Command     cmd.Command
}

type process struct {
	dir     string
	wait    *sync.WaitGroup
	command cmd.Command
}

func buildThreeComponentPath(a string, b string, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func (bp *BuildProcess) Do() {
	start := time.Now()
	sem := make(chan struct{}, bp.ProcessParams.Parallel)
	processChannels := make(chan process)
	rootDir := bp.ProcessParams.Root
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
			filePath := buildThreeComponentPath(rootDir, name, bp.ProcessParams.CheckFile)
			command := bp.ProcessParams.Command
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println("add building process -" + name)
				waitGroup.Add(1)
				fp := rootDir + "/" + name
				cutExec := bp.ProcessParams.CutExecPath
				replacedFullPath := strings.ReplaceAll(fp, cutExec, "")
				processChannels <- process{
					replacedFullPath,
					&waitGroup,
					command,
				}
			}
			if len(bp.ProcessParams.Recursive) > 0 {
				for _, recursiveDir := range bp.ProcessParams.Recursive {
					recursive := buildThreeComponentPath(rootDir, name, recursiveDir)
					recursive = strings.ReplaceAll(recursive, recursiveDir+"/"+recursiveDir, recursiveDir)
					if _, err := os.Stat(recursive); err == nil {
						newBp := BuildProcess{
							Mode: bp.Mode,
							ProcessParams: &ProcessParams{
								Root:        recursive,
								CutExecPath: bp.ProcessParams.CutExecPath,
								Parallel:    bp.ProcessParams.Parallel,
								Recursive:   bp.ProcessParams.Recursive,
								CheckFile:   bp.ProcessParams.CheckFile,
								Command:     command,
							},
						}
						newBp.Do()
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
	log.Printf("Elapsed time for %s - %s", rootDir, elapsed)
}

func handlePrc(prc process, sem chan struct{}) error {
	sem <- struct{}{}
	defer func() {
		<-sem
		prc.wait.Done()
	}()
	prc.command.AddArg(prc.dir)
	return prc.command.Run()
}
