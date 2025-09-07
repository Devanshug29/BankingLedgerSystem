package config

type Config struct {
	PostgresURL string `yaml:"postgres_url"`

	MongoURL string `yaml:"mongo_url"`
	MongoDB  string `yaml:"mongo_db"`

	APIPort          string `yaml:"api_port"`
	ProcessorWorkers int    `yaml:"processor_workers"`
	
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type KafkaConfig struct {
	Brokers     []string `mapstructure:"brokers"`
	Topic       string   `mapstructure:"topic"`
	ClientID    string   `mapstructure:"clientID"`
	Acks        string   `mapstructure:"acks"`
	Retries     int      `mapstructure:"retries"`
	Partitioner string   `mapstructure:"partitioner"`
}
