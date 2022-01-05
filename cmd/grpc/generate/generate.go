package generate

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"task-runner-cobra/cmd/custom"
	"task-runner-cobra/utils/command/custom/out"
)

const (
	ServiceTypeClient = "client"
	ServiceTypeServer = "server"
)

type generationProcess struct {
	ServiceName  string
	ServiceType  string
	ProtocParams *protocParams
}

type protocParams struct {
	Root   string
	Plugin string
	Out    string
	Common string
	Clear  []string
}

func init() {
	Cmd.PersistentFlags().String("pattern", "", "example: [client|server]:[service name]")
	_ = Cmd.MarkPersistentFlagRequired("pattern")
}

var (
	Cmd = &cobra.Command{
		Use:   "grpc-generate-proto",
		Short: "Generate proto files",
		Long: `
Generate airo proto files (client or server) from plugin.
Available variables in commands:
	<-root> - root from config file
	<-plugin> - plugin from config file
	<-out> - out from config file
	<-files> - found proto files
`,
		Run: func(cmd *cobra.Command, args []string) {
			pattern, _ := cmd.Flags().GetString("pattern")
			typeAndService := strings.Split(pattern, ":")
			serviceType := typeAndService[0]
			if serviceType != ServiceTypeClient && serviceType != ServiceTypeServer {
				panic("Unknown compilation type")
			}
			var config interface{}
			config = viper.Get("grpc." + serviceType)
			var serviceName string
			if cap(typeAndService) == 2 {
				serviceName = typeAndService[1]
			}

			configMap := config.(map[string]interface{})
			pp := &protocParams{
				configMap["root"].(string),
				configMap["plugin"].(string),
				configMap["out"].(string),
				configMap["common"].(string),
				viper.GetStringSlice("grpc." + serviceType + ".clear"),
			}

			genProcess := generationProcess{
				serviceName,
				serviceType,
				pp,
			}
			pathToService := genProcess.ProtocParams.Root + "/" + genProcess.ServiceName
			if len(genProcess.ProtocParams.Clear) > 0 {
				for _, dir := range genProcess.ProtocParams.Clear {
					err := os.RemoveAll(dir)
					if err != nil {
						panic(err)
					}
				}
			}
			if _, err := os.Stat(genProcess.ProtocParams.Out); err != nil {
				err := os.Mkdir(genProcess.ProtocParams.Out, os.ModePerm)
				if err != nil {
					panic(err)
				}
			}
			proto := genProcess.findProto(pathToService)
			if genProcess.ProtocParams.Common != "" {
				proto = append(proto, genProcess.findProto(genProcess.ProtocParams.Common)...)
			}
			genProcess.runProtoc(proto)
		},
	}
)

func (gp *generationProcess) findProto(root string) []string {
	var proto []string
	if _, err := os.Stat(root); err != nil {
		return proto
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), "imports.proto") {
			return nil
		}
		if !info.IsDir() && filepath.Ext(strings.TrimSpace(path)) == ".proto" {
			proto = append(proto, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return proto
}

func (gp *generationProcess) runProtoc(proto []string) {
	customCommands := custom.BuildCommandsWithViaConfig("grpc." + gp.ServiceType + ".commands")
	if len(customCommands) == 0 {
		args := []string{
			"--proto_path=<-root>",
			"--plugin=protoc-gen-grpc=<-plugin>",
			"--php_out=<-out>",
			"--grpc_out=<-out>",
			"<-files>",
		}
		command := &custom.Command{
			Main: "protoc",
			Args: args,
		}
		customCommands = append(customCommands, command)
	}
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	for _, customCommand := range customCommands {
		customCommand.ReplaceArgs(map[string]string{
			"<-root>":   gp.ProtocParams.Root,
			"<-plugin>": gp.ProtocParams.Plugin,
			"<-out>":    gp.ProtocParams.Out,
		})
		if customCommand.ArgsContains("<-files>") {
			customCommand.RemoveArg(customCommand.FindArgIndex("<-files>"))
			for _, filePath := range proto {
				fmt.Println(green("•"), bold(filePath))
				customCommand.AddArg(filePath)
			}
		}
	}
	out.HandleBatch(customCommands)
}
