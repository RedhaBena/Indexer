package config

var (
	DefaultLogLevel string = "debug"
)

type LoggerConfig struct {
	LogLevel string
}
