package config

type Config struct {
	PostgresURL string `yaml:"postgres_url"`

	MongoURL string `yaml:"mongo_url"`
	MongoDB  string `yaml:"mongo_db"`

	KafkaBrokers string `yaml:"kafka_brokers"`
	KafkaTopic   string `yaml:"kafka_topic"`

	APIPort          string `yaml:"api_port"`
	ProcessorWorkers int    `yaml:"processor_workers"`
}
