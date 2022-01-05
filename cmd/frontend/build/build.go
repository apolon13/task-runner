package build

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"task-runner-cobra/cmd/custom"
	commandOut "task-runner-cobra/utils/command/custom/out"
	"time"
)

type buildProcess struct {
	Mode               string
	BuildProcessParams *buildProcessParams
}

type buildProcessParams struct {
	Root        string
	CutExecPath string
	Parallel    int
	CheckFile   string
}

type process struct {
	wait     *sync.WaitGroup
	commands []*custom.Command
}

func init() {
	Cmd.PersistentFlags().String("mode", "production", "build mode")
}

var (
	Cmd = &cobra.Command{
		Use:   "build-frontend",
		Short: "Vue frontend builder",
		Long: `
Parallel assembly of vue modules.
Available variables in commands:
	<-mode> - build mode
	<-module> - current module in process
`,
		Run: func(cmd *cobra.Command, args []string) {
			mode, _ := cmd.Flags().GetString("mode")

			bp := &buildProcess{
				Mode: mode,
				BuildProcessParams: &buildProcessParams{
					viper.GetString("build.frontend.root"),
					viper.GetString("build.frontend.cut-exec-path"),
					viper.GetInt("build.frontend.parallel"),
					viper.GetString("build.frontend.check-file"),
				},
			}

			start := time.Now()
			sem := make(chan struct{}, bp.BuildProcessParams.Parallel)
			processChannels := make(chan process)
			rootDir := bp.BuildProcessParams.Root
			var modules []string
			fmt.Println("root directory - " + rootDir)
			err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
				if strings.Contains(path, bp.BuildProcessParams.CheckFile) && err == nil {
					path = strings.ReplaceAll(path, "/"+bp.BuildProcessParams.CheckFile, "")
					modules = append(modules, strings.ReplaceAll(path, bp.BuildProcessParams.CutExecPath, ""))
				}
				return err
			})
			if err != nil {
				panic(err)
			}
			go func() {
				var waitGroup sync.WaitGroup
				for _, module := range modules {
					fmt.Println("add building process -" + module)
					waitGroup.Add(1)
					customCommands := custom.BuildCommandsWithViaConfig("build.frontend.commands")
					for _, customCommand := range customCommands {
						customCommand.ReplaceArgs(map[string]string{
							"<-mode>":   mode,
							"<-module>": module,
						})
					}
					processChannels <- process{
						&waitGroup,
						customCommands,
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
			fmt.Printf("Elapsed time for %s - %s\n", rootDir, elapsed)
		},
	}
)

func handlePrc(prc process, sem chan struct{}) {
	sem <- struct{}{}
	defer func() {
		<-sem
		prc.wait.Done()
	}()
	commandOut.HandleBatch(prc.commands)
}
