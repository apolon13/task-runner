package db

import (
	"fmt"
	"github.com/fatih/color"
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
	stdErr := make(chan string)
	stdOut := make(chan string)
	quit := make(chan struct{})
	go func() {
		err := r.Command.Run(stdErr, stdOut)
		if err != nil {
			log.Fatal(err)
		}
		quit <- struct{}{}
	}()
	for {
		select {
		case errString := <-stdErr:
			color.Red(errString)
		case outString := <-stdOut:
			color.Blue(outString)
		case <-quit:
			return
		}
	}
}
