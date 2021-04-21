package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"task-runner/config"
)

func Handle(cnfCommand config.Command) error {
	cmd := exec.Command(cnfCommand.Main, cnfCommand.Args...)
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
	return cmd.Run()
}
