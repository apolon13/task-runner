package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"task-runner/config"
	s3Connection "task-runner/connection/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

type DownloadFile struct {
	LocalPath  string
	RemotePath string
	FileName   string
	Client     *s3Connection.Client
}

func (df *DownloadFile) GetFileName() string {
	ss := strings.Split(df.FileName, "/")
	s := ss[len(ss)-1]
	return s
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percentageDownloaded := float32(pw.written*100) / float32(pw.size)

	fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\r", pw.size, pw.written, percentageDownloaded)

	return pw.writer.WriteAt(p, off)
}

func getFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func NewFile(cnf config.Yaml, filename string, client *s3Connection.Client) *DownloadFile {
	return &DownloadFile{
		LocalPath:  cnf.Restore.Db.Path.S3.Local,
		RemotePath: cnf.Restore.Db.Path.S3.Remote,
		FileName:   filename,
		Client:     client,
	}
}

func (df *DownloadFile) localFullPath() string {
	return fmt.Sprintf("%s/%s", df.LocalPath, df.GetFileName())
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

	bucket := df.RemotePath
	key := df.FileName

	downloader := s3manager.NewDownloader(df.Client.Session)
	size, err := getFileSize(df.Client.Connection, bucket, key)
	if err != nil {
		panic(err)
	}

	temp, err := ioutil.TempFile(df.LocalPath, "getObjWithProgress-tmp-")
	if err != nil {
		panic(err)
	}
	tempFileName := temp.Name()
	writer := &progressWriter{writer: temp, size: size, written: 0}
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err := downloader.Download(writer, params); err != nil {
		err := os.Remove(tempFileName)
		if err != nil {
			panic(err)
		}
	}

	if err := temp.Close(); err != nil {
		panic(err)
	}

	if err := os.Rename(tempFileName, df.localFullPath()); err != nil {
		panic(err)
	}
}
