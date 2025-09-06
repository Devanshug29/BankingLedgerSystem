package config

type Config struct {
	PostgresHost     string `yaml:"postgres_host"`
	PostgresPort     string `yaml:"postgres_port"`
	PostgresDB       string `yaml:"postgres_db"`
	PostgresUser     string `yaml:"postgres_user"`
	PostgresPassword string `yaml:"postgres_password"`
	PostgresSSLMode  string `yaml:"postgres_ssl_mode"`

	MongoURI string `yaml:"mongo_uri"`
	MongoDB  string `yaml:"mongo_db"`

	KafkaBrokers string `yaml:"kafka_brokers"`
	KafkaTopic   string `yaml:"kafka_topic"`

	APIPort          string `yaml:"api_port"`
	ProcessorWorkers int    `yaml:"processor_workers"`
}
