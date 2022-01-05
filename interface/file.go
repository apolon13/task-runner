package file

type DownloadFile interface {
	Process()
	RemoveLocal() error
	GetFileName() string
}
