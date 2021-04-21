package backup

import (
	"fmt"
	"log"
	"strings"
	"task-runner/cmd"
	"task-runner/config"
	"task-runner/downloader/sftp"
)

func Do(df *sftp.DownloadFile, cnf config.Yaml, db string) error {
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
	cnf.Backup.Command.Args = args
	return cmd.Handle(cnf.Backup.Command)
}
