package db

import (
	"fmt"
	"log"
	"task-runner/cmd"
	"task-runner/downloader/sftp"
)

type Restore struct {
	File    *sftp.DownloadFile
	Command *cmd.Command
	Remove  bool
}

func (r *Restore) Do() {
	fmt.Println(fmt.Sprintf("Download file - %s", r.File.FileName))
	r.File.Process()
	fmt.Println("Downloading complete")
	if r.Remove == true {
		defer func() {
			if err := r.File.RemoveLocal(); err != nil {
				log.Fatal(fmt.Errorf("error removing downloaded file: %s", err))
			}
		}()
	}
	err := r.Command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
