package frontend

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
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
	Root        string `yaml:"root"`
	CutExecPath string `yaml:"cut-exec-path"`
	Parallel    int    `yaml:"parallel"`
	CheckFile   string `yaml:"check-file"`
	Command     cmd.Command
}

type process struct {
	dir     string
	wait    *sync.WaitGroup
	command cmd.Command
	mode    string
}

func (bp *BuildProcess) Do() {
	start := time.Now()
	sem := make(chan struct{}, bp.ProcessParams.Parallel)
	processChannels := make(chan process)
	rootDir := bp.ProcessParams.Root
	var modules []string
	fmt.Println("root directory - " + rootDir)
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, bp.ProcessParams.CheckFile) && err == nil {
			path = strings.ReplaceAll(path, "/"+bp.ProcessParams.CheckFile, "")
			modules = append(modules, strings.ReplaceAll(path, bp.ProcessParams.CutExecPath, ""))
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	go func() {
		var waitGroup sync.WaitGroup
		for _, module := range modules {
			command := bp.ProcessParams.Command
			fmt.Println("add building process -" + module)
			waitGroup.Add(1)
			processChannels <- process{
				module,
				&waitGroup,
				command,
				bp.Mode,
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

func handlePrc(prc process, sem chan struct{}) {
	sem <- struct{}{}
	defer func() {
		<-sem
		prc.wait.Done()
	}()
	prc.command.AddArg(prc.dir)
	stdErr := make(chan string)
	stdOut := make(chan string)
	quit := make(chan struct{})
	go func() {
		err := prc.command.Run(stdErr, stdOut)
		if err != nil {
			log.Fatal(err)
		}
		quit <- struct{}{}
	}()
	for {
		select {
		case errString := <-stdErr:
			switch prc.mode {
			case "production":
				color.Red(errString)
				log.Fatal(fmt.Sprintf("Error in build %s", prc.command))
			case "development":
				color.Red(errString)
			}
		case outString := <-stdOut:
			color.Blue(outString)
		case <-quit:
			return
		}
	}
}
