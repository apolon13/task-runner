package frontend

import (
	"backup-downloader/config"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type BuildParams struct {
	Cnf  config.Yaml
	Mode string
}

type command struct {
	main string
	args []string
}

type buildProcess struct {
	dir     string
	mode    string
	wait    *sync.WaitGroup
	command command
}

func buildThreeComponentPath(a string, b string, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func Do(bp BuildParams) {
	start := time.Now()
	sem := make(chan struct{}, bp.Cnf.Build.Frontend.Parallel)
	processChannels := make(chan buildProcess)
	fmt.Println("root directory - " + bp.Cnf.Build.Frontend.Root)
	files, err := ioutil.ReadDir(bp.Cnf.Build.Frontend.Root)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		var waitGroup sync.WaitGroup
		for _, file := range files {
			if file.IsDir() {
				filePath := buildThreeComponentPath(bp.Cnf.Build.Frontend.Root, file.Name(), bp.Cnf.Build.Frontend.IfExistFile)
				if _, err := os.Stat(filePath); err == nil {
					fmt.Println("add building process -" + file.Name())
					waitGroup.Add(1)
					fullPath := bp.Cnf.Build.Frontend.Root + "/" + file.Name()
					replacedFullPath := strings.ReplaceAll(fullPath, bp.Cnf.Build.Frontend.ClearPath, "")
					processChannels <- buildProcess{
						replacedFullPath,
						bp.Mode,
						&waitGroup,
						command{
							main: bp.Cnf.Build.Frontend.Command.Main,
							args: bp.Cnf.Build.Frontend.Command.Args,
						},
					}
				}
				if len(bp.Cnf.Build.Frontend.Recursive) > 0 {
					for _, recursiveDir := range bp.Cnf.Build.Frontend.Recursive {
						recursive := buildThreeComponentPath(bp.Cnf.Build.Frontend.Root, file.Name(), recursiveDir)
						recursive = strings.ReplaceAll(recursive, recursiveDir + "/" + recursiveDir, recursiveDir)
						if _, err := os.Stat(recursive); err == nil {
							cnf := bp.Cnf
							cnf.Build.Frontend.Root = recursive
							Do(BuildParams{
								Cnf: cnf,
								Mode: bp.Mode,
							})
						}
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

func handlePrc(prc buildProcess, sem chan struct{}) {
	sem <- struct{}{}
	defer func() {
		<-sem
		prc.wait.Done()
	}()
	var args []string
	for _, arg := range prc.command.args {
		arg = strings.ReplaceAll(arg, "${-mode}", prc.mode)
		args = append(args, arg)
	}
	args = append(args, prc.dir)
	cmd := exec.Command(prc.command.main, args...)
	if cmd.Stderr == nil {
		cmdErrReader, err := cmd.StderrPipe()
		errScanner := bufio.NewScanner(cmdErrReader)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for errScanner.Scan() {
				fmt.Println("error: " + errScanner.Text())
			}
		}()
	}

	if cmd.Stdout == nil {
		cmdOutReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		outScanner := bufio.NewScanner(cmdOutReader)
		go func() {
			for outScanner.Scan() {
				fmt.Println(outScanner.Text())
			}
		}()
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
