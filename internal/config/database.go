package config

var (
	DefaultDatabaseHost string = "localhost:7687"
	DefaultDatabaseUser string = "neo4j"
	DefaultDatabasePass string = "aztec-peace-linear-laura-gregory-4537"
)

type DatabaseConfig struct {
	Host string
	User string
	Pass string
}
