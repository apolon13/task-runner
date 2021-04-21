package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

type DownloadFile struct {
	LocalPath          string
	RemotePath string
	FileName               string
	Connection             *ssh.Client
}

func (df *DownloadFile) localFullPath() string {
	return fmt.Sprintf("%s/%s", df.LocalPath, df.FileName)
}

func (df *DownloadFile) remoteFullPath() string {
	return fmt.Sprintf("%s/%s", df.RemotePath, df.FileName)
}

func (df *DownloadFile) RemoveLocal() error {
	return os.Remove(df.localFullPath())
}

func (df *DownloadFile) Process() {
	if _, err := os.Stat(df.LocalPath); os.IsNotExist(err) {
		if err = os.Mkdir(df.LocalPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	sftpWrap, err := sftp.NewClient(df.Connection)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpWrap.Close()
	localFile, err := os.Create(df.localFullPath())
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()

	remoteFile, err := sftpWrap.Open(df.remoteFullPath())
	if err != nil {
		log.Fatal(err)
	}
	defer remoteFile.Close()

	if _, err := remoteFile.WriteTo(localFile); err != nil {
		panic(err)
	}
}
