package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"task-runner/config"
	"task-runner/connection/ssh"
	"task-runner/downloader/sftp"
	"task-runner/utils/backup"
	"task-runner/utils/builder/frontend"
	gitUtil "task-runner/utils/git"
)

func buildConfig(path string) config.Yaml {
	yamlFile := config.Yaml{}
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("config file not found: %s", err))
	}
	err = yaml.Unmarshal(configFile, &yamlFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("unmarshal config file error: %s", err))
	}
	return yamlFile
}

func main() {
	wd, _ := os.Getwd()
	defaultConfigFile := wd + "/config.yaml"
	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	backupCnf := backupCmd.String("cnf", defaultConfigFile, "config file path")
	f := backupCmd.String("f", "", "backup file name")
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

	switch os.Args[1] {
	case "backup":
		_ = backupCmd.Parse(os.Args[2:])
		yamlFile := buildConfig(*backupCnf)
		client := &ssh.Client{
			Params: ssh.Params{
				Username:   yamlFile.Connections.Ssh.Username,
				Host:       yamlFile.Connections.Ssh.Host,
				Port:       yamlFile.Connections.Ssh.Port,
				PrivateKey: yamlFile.Connections.Ssh.PrivateKey,
				Password:   yamlFile.Connections.Ssh.Password,
			},
		}
		client.Connect()
		defer client.Connection.Close()
		df := &sftp.DownloadFile{
			LocalPath:  yamlFile.Backup.Path.Local,
			RemotePath: yamlFile.Backup.Path.Remote,
			FileName:   *f,
			Connection: client.Connection,
		}
		err := backup.Do(df, yamlFile, *db)
		if err != nil {
			log.Fatal(err)
		}
	case "build-frontend":
		_ = buildFrontendCmd.Parse(os.Args[2:])
		yamlFile := buildConfig(*buildFrontendCnf)
		frontend.Do(frontend.BuildParams{
			Cnf:  yamlFile,
			Mode: *mode,
		})
	case "release":
		_ = releaseCmd.Parse(os.Args[2:])
		branch := *releaseBranch
		if branch == "current" {
			branch, _ = gitUtil.CurrentBranch()
		}
		gitUtil.Release(
			&config.Branch{Name: strings.Trim(branch, "\n")},
			buildConfig(*releaseCnf))
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
		os.Exit(2)
	default:
		fmt.Println("Undefined subcommand")
		os.Exit(1)
	}
}
