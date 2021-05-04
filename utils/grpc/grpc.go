package grpc

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"task-runner/cmd"
)

const (
	CompilationTypeClient = "client"
	CompilationTypeServer = "server"
)

type CompilationProcess struct {
	ServiceName  string
	ProtocParams ProtocParams
}

type ProtocParams struct {
	Root   string   `yaml:"root"`
	Plugin string   `yaml:"plugin"`
	Out    string   `yaml:"out"`
	Common string   `yaml:"common"`
	Clear  []string `yaml:"clear"`
}

func (cp *CompilationProcess) Do() {
	pathToService := cp.ProtocParams.Root + "/" + cp.ServiceName
	if len(cp.ProtocParams.Clear) > 0 {
		for _, dir := range cp.ProtocParams.Clear {
			err := os.RemoveAll(dir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if _, err := os.Stat(cp.ProtocParams.Out); err != nil {
		err := os.Mkdir(cp.ProtocParams.Out, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	proto := cp.findProto(pathToService)
	if cp.ProtocParams.Common != "" {
		proto = append(proto, cp.findProto(cp.ProtocParams.Common)...)
	}
	cp.runProtoc(proto)
}

func (cp *CompilationProcess) findProto(root string) []string {
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
		log.Fatal(err)
	}
	return proto
}

func (cp *CompilationProcess) runProtoc(proto []string) {
	args := []string{
		"--proto_path=" + cp.ProtocParams.Root,
		"--plugin=protoc-gen-grpc=" + cp.ProtocParams.Plugin,
		"--php_out=" + cp.ProtocParams.Out,
		"--grpc_out=" + cp.ProtocParams.Out,
	}
	command := &cmd.Command{
		Main: "protoc",
		Args: append(args, proto...),
	}
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	for _, filePath := range proto {
		fmt.Println(green("•"), bold(filePath))
	}
}
