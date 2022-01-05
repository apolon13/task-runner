package custom

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
	envFinder "task-runner-cobra/utils/env/finder"
)

type Command struct {
	Main         string
	Args         []string
	IgnoreStdErr bool
}

func (c *Command) prepare() {
	for _, arg := range c.Args {
		envMap := envFinder.FindEnvBetweenQuotes(arg)
		c.ReplaceArgs(envMap)
	}
	c.Main = replaceAllTokensInString(c.Main, envFinder.FindEnvBetweenQuotes(c.Main))
}

func replaceAllTokensInString(str string, tokensAndValues map[string]string) string {
	for token, value := range tokensAndValues {
		str = strings.ReplaceAll(str, token, value)
	}
	return str
}

func (c *Command) Run(stderr chan<- string, stdout chan<- string) error {
	c.prepare()
	cmd := exec.Command(c.Main, c.Args...)
	if cmd.Stderr == nil {
		cmdErrReader, err := cmd.StderrPipe()
		errScanner := bufio.NewScanner(cmdErrReader)
		if err != nil {
			panic(err)
		}
		go func() {
			for errScanner.Scan() {
				stderr <- errScanner.Text()
			}
		}()
	}

	if cmd.Stdout == nil {
		cmdOutReader, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		outScanner := bufio.NewScanner(cmdOutReader)
		go func() {
			for outScanner.Scan() {
				stdout <- outScanner.Text()
			}
		}()
	}
	color.Blue(cmd.String())
	return cmd.Run()
}

func (c *Command) ReplaceArgs(args map[string]string) {
	var newArgs []string
	for _, arg := range c.Args {
		newArgs = append(newArgs, replaceAllTokensInString(arg, args))
	}
	c.Args = newArgs
}

func (c *Command) AddArg(arg string) {
	c.Args = append(c.Args, arg)
}

func (c *Command) FindArgIndex(arg string) int {
	for index, val := range c.Args {
		if val == arg {
			return index
		}
	}
	return -1
}

func (c *Command) RemoveArg(index int) {
	c.Args = append(c.Args[:index], c.Args[index+1:]...)
}

func (c *Command) ArgsContains(arg string) bool {
	for _, v := range c.Args {
		if v == arg {
			return true
		}
	}
	return false
}

func BuildCommandsWithViaConfig(path string) []*Command {
	commands := viper.Get(path)
	var customCommands []*Command
	if commands != nil {
		for _, command := range commands.([]interface{}) {
			commandMap := command.(map[interface{}]interface{})
			var commandArgs []string
			for _, arg := range commandMap["args"].([]interface{}) {
				commandArgs = append(commandArgs, arg.(string))
			}
			ignoreStdErr := commandMap["ignore-std-err"]
			if ignoreStdErr == nil {
				ignoreStdErr = false
			}
			customCommand := Command{
				Main:         commandMap["main"].(string),
				Args:         commandArgs,
				IgnoreStdErr: ignoreStdErr.(bool),
			}

			customCommands = append(customCommands, &customCommand)
		}
	}
	return customCommands
}
