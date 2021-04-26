package cmd

import (
	"bufio"
	"github.com/fatih/color"
	"log"
	"os/exec"
	"strings"
)

type Command struct {
	Main string   `yaml:"main"`
	Args []string `yaml:"args"`
}

func (c *Command) Run() error {
	cmd := exec.Command(c.Main, c.Args...)
	if cmd.Stderr == nil {
		cmdErrReader, err := cmd.StderrPipe()
		errScanner := bufio.NewScanner(cmdErrReader)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for errScanner.Scan() {
				color.Red(errScanner.Text())
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
				color.Blue(outScanner.Text())
			}
		}()
	}
	color.Green(cmd.String())
	return cmd.Run()
}

func (c *Command) ReplaceArgs(args map[string]string) {
	var newArgs []string
	for _, arg := range c.Args {
		for token, value := range args {
			arg = strings.ReplaceAll(arg, token, value)
		}
		newArgs = append(newArgs, arg)
	}
	c.Args = newArgs
}

func (c *Command) AddArg(arg string) {
	c.Args = append(c.Args, arg)
}
