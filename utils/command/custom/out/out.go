package out

import (
	"github.com/fatih/color"
	"task-runner-cobra/cmd/custom"
)

func Handle(command *custom.Command) {
	stdErr := make(chan string)
	stdOut := make(chan string)
	quit := make(chan struct{})
	go func() {
		err := command.Run(stdErr, stdOut)
		if err != nil {
			stdErr <- err.Error()
		}
		quit <- struct{}{}
	}()
	for {
		select {
		case errString := <-stdErr:
			if command.IgnoreStdErr == false {
				color.Red(errString)
				panic(errString)
			}
		case outString := <-stdOut:
			color.Green(outString)
		case <-quit:
			return
		}
	}
}

func HandleBatch(commands []*custom.Command) {
	for _, command := range commands {
		Handle(command)
	}
}
