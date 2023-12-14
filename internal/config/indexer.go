package config

var (
	DefaultBatchSize uint = 2000
)

type IndexerConfig struct {
	BatchSize uint
}
