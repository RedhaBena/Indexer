package config

var (
	DefaultFilePath string = "biggertest.json"
)

type FileConfig struct {
	LocalPath    string
	DownloadPath string
}
