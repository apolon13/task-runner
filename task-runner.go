package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"task-runner/config"
	s3Connection "task-runner/connection/s3"
	sshConnection "task-runner/connection/ssh"
	"task-runner/downloader/file"
	s3Downloader "task-runner/downloader/s3"
	sftpDownloader "task-runner/downloader/sftp"
	"task-runner/utils/builder/frontend"
	"task-runner/utils/grpc"
	dbUtil "task-runner/utils/restore/db"
	"task-runner/utils/services"
)

func main() {
	wd, _ := os.Getwd()
	defaultConfigFile := wd + "/config.yaml"
	restoreCmd := flag.NewFlagSet("restore-db", flag.ExitOnError)
	restoreCnf := restoreCmd.String("cnf", defaultConfigFile, "config file path")
	f := restoreCmd.String("f", "", "restore file name")
	db := restoreCmd.String("db", "", "database")
	con := restoreCmd.String("con", "ssh", "connection name")

	buildFrontendCmd := flag.NewFlagSet("build-frontend", flag.ExitOnError)
	buildFrontendCnf := buildFrontendCmd.String("cnf", defaultConfigFile, "config file path")
	mode := buildFrontendCmd.String("mode", "production", "production or development")

	grpcCmd := flag.NewFlagSet("grpc", flag.ExitOnError)
	pattern := grpcCmd.String("pattern", "", "<client or server>[:<service_name>]")
	grpcCnf := grpcCmd.String("cnf", defaultConfigFile, "config file path")

	servicesInfoCmd := flag.NewFlagSet("services-info", flag.ExitOnError)
	servicesInfoCnf := servicesInfoCmd.String("cnf", defaultConfigFile, "config file path")
	servicesInfoFile := servicesInfoCmd.String("f", "", "export to file")

	switch os.Args[1] {
	case "restore-db":
		_ = restoreCmd.Parse(os.Args[2:])
		yamlFile := config.New(*restoreCnf)
		var df file.DownloadFile
		switch *con {
		case "ssh":
			df = sftpDownloader.NewFile(yamlFile, *f, sshConnection.NewClient(yamlFile))
		case "s3":
			df = s3Downloader.NewFile(yamlFile, *f, s3Connection.NewClient(yamlFile))
		default:
			fmt.Println("Undefined connection")
			os.Exit(1)
		}

		command := &yamlFile.Restore.Db.Command
		command.ReplaceArgs(map[string]string{
			"<-f>":  df.GetFileName(),
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
	case "services-info":
		yamlFile := config.New(*servicesInfoCnf)
		_ = servicesInfoCmd.Parse(os.Args[2:])
		services.PrintInfo(sshConnection.NewClient(yamlFile), *servicesInfoFile)
	case "-h":
		fmt.Println("Usage: task-runner " + restoreCmd.Name())
		restoreCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + buildFrontendCmd.Name())
		buildFrontendCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + grpcCmd.Name())
		grpcCmd.PrintDefaults()
		fmt.Println("Usage: task-runner " + servicesInfoCmd.Name())
		servicesInfoCmd.PrintDefaults()
		os.Exit(2)
	default:
		fmt.Println("Undefined subcommand")
		os.Exit(1)
	}
}
