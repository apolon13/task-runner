package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"task-runner/config"
	"task-runner/connection/ssh"
	"task-runner/downloader/sftp"
	"task-runner/utils/builder/frontend"
	gitUtil "task-runner/utils/git"
	"task-runner/utils/grpc"
	dbUtil "task-runner/utils/restore/db"
)

func main() {
	wd, _ := os.Getwd()
	defaultConfigFile := wd + "/config.yaml"
	backupCmd := flag.NewFlagSet("restore-db", flag.ExitOnError)
	backupCnf := backupCmd.String("cnf", defaultConfigFile, "config file path")
	f := backupCmd.String("f", "", "restore file name")
	db := backupCmd.String("db", "", "database")

	buildFrontendCmd := flag.NewFlagSet("build-frontend", flag.ExitOnError)
	buildFrontendCnf := buildFrontendCmd.String("cnf", defaultConfigFile, "config file path")
	mode := buildFrontendCmd.String("mode", "production", "production or development")

	releaseCmd := flag.NewFlagSet("release", flag.ExitOnError)
	releaseBranch := releaseCmd.String("branch", "current", "release branch")
	releaseCnf := releaseCmd.String("cnf", defaultConfigFile, "config file path")

	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	testStand := deployCmd.String("stand", "", "test stand")
	deployBranch := deployCmd.String("branch", "current", "deploy branch")

	grpcCmd := flag.NewFlagSet("grpc", flag.ExitOnError)
	pattern := grpcCmd.String("pattern", "", "<client or server>[:<service_name>]")
	grpcCnf := grpcCmd.String("cnf", defaultConfigFile, "config file path")

	switch os.Args[1] {
	case "restore-db":
		_ = backupCmd.Parse(os.Args[2:])
		yamlFile := config.New(*backupCnf)
		df := sftp.NewFile(yamlFile, *f, ssh.NewClient(yamlFile))
		command := &yamlFile.Restore.Db.Command
		command.ReplaceArgs(map[string]string{
			"<-f>":  df.FileName,
			"<-db>": *db,
		})
		restore := &dbUtil.Restore{
			File:    df,
			Command: command,
			Remove:  yamlFile.Restore.Db.Remove,
		}
		restore.Do()
	case "build-frontend":
		_ = buildFrontendCmd.Parse(os.Args[2:])
		yamlFile := config.New(*buildFrontendCnf)
		params := &yamlFile.Build.Frontend
		params.Command.ReplaceArgs(map[string]string{
			"<-mode>": *mode,
		})
		bp := &frontend.BuildProcess{
			Mode:          *mode,
			ProcessParams: params,
		}
		bp.Do()
	case "grpc":
		_ = grpcCmd.Parse(os.Args[2:])
		yamlFile := config.New(*grpcCnf)
		if *pattern == "" {
			log.Fatal("Pattern missing")
		}
		typeAndService := strings.Split(*pattern, ":")
		var params grpc.ProtocParams
		switch typeAndService[0] {
		case grpc.CompilationTypeClient:
			params = yamlFile.GRPC.Client
		case grpc.CompilationTypeServer:
			params = yamlFile.GRPC.Server
		default:
			log.Fatal("Unknown compilation type")
		}
		var serviceName string
		if cap(typeAndService) == 2 {
			serviceName = typeAndService[1]
		}
		cp := &grpc.CompilationProcess{
			ServiceName:  serviceName,
			ProtocParams: params,
		}
		cp.Do()
	case "release":
		_ = releaseCmd.Parse(os.Args[2:])
		branch := *releaseBranch
		if branch == "current" {
			branch, _ = gitUtil.CurrentBranch()
		}
		gitUtil.Release(
			&config.Branch{Name: strings.Trim(branch, "\n")},
			config.New(*releaseCnf))
	case "deploy":
		_ = deployCmd.Parse(os.Args[2:])
		branch := *deployBranch
		if branch == "current" {
			branch, _ = gitUtil.CurrentBranch()
		}
		if *testStand == "" {
			log.Fatal("Test stand missing")
		}
		gitUtil.Deploy(
			&config.Branch{Name: strings.Trim(branch, "\n")},
			*testStand)
	case "-h":
		fmt.Println("Usage: task-runner " + backupCmd.Name())
		backupCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + buildFrontendCmd.Name())
		buildFrontendCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + releaseCmd.Name())
		releaseCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + deployCmd.Name())
		deployCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + grpcCmd.Name())
		grpcCmd.PrintDefaults()
		os.Exit(2)
	default:
		fmt.Println("Undefined subcommand")
		os.Exit(1)
	}
}
