package backup

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"task-runner/config"
	"task-runner/downloader/sftp"
)

func Do(df *sftp.DownloadFile, cnf config.Yaml, db string) {
	fmt.Println(fmt.Sprintf("Download file - %s", df.FileName))
	df.Process()
	fmt.Println("Downloading complete")
	if cnf.Backup.Remove == true {
		defer func() {
			if err := df.RemoveLocal(); err != nil {
				log.Fatal(fmt.Errorf("error removing downloaded file: %s", err))
			}
		}()
	}
	var args []string
	for _, arg := range cnf.Backup.Command.Args {
		arg = strings.ReplaceAll(arg, "${-f}", df.FileName)
		arg = strings.ReplaceAll(arg, "${-db}", db)
		args = append(args, arg)
	}
	cmd := exec.Command(cnf.Backup.Command.Main, args...)
	cmdStdOutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdStdErrorPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("Executing command - %s", cmd.String()))
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	io.Copy(os.Stdout, cmdStdOutPipe)
	io.Copy(os.Stderr, cmdStdErrorPipe)
}
