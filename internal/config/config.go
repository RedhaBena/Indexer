package config

var GlobalConfig Config = Config{}

type Config struct {
	DatabaseConfig
	FileConfig
	IndexerConfig
	LoggerConfig
}
