package config

// MQConfig defines the configuration for the message queue.
type MQConfig struct {
	Dialect     string `envconfig:"MQ_DIALECT"`
	Debug       bool   `envconfig:"MQ_DEBUG"`
	Track       bool   `envconfig:"MQ_TRACK"`
	DeadLetterQ string `envconfig:"DEAD_LETTER_QUEUE"`
	HandlerName string `envconfig:"HANDLER_NAME"`
	Redis       RedisConfig
	Amqp        AmqpConfig
	Kafka       KafkaConfig
	GoogleCloud GoogleCloud
	SQL         SQL
}

// RedisConfig defines the configuration for Redis.
type RedisConfig struct {
	RedisURL      string `envconfig:"REDIS_URI"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP"`
	UserName      string `envconfig:"REDIS_USERNAME"`
	Password      string `envconfig:"REDIS_PASSWORD"`
}

// AmqpConfig defines the configuration for AMQP.
type AmqpConfig struct {
	AmqpURL string `envconfig:"AMQB_URI"`
}

// KafkaConfig defines the configuration for Kafka.
type KafkaConfig struct {
	KafkaBroker   []string `envconfig:"KAFKA_BROKER"`
	ConsumerGroup string   `envconfig:"CONSUMER_GROUP"`
}

// GoogleCloud defines the configuration for Google Cloud Pub/Sub.
type GoogleCloud struct {
	ProjectID      string `envconfig:"PROJECT_ID"`
	SubscriptionID string `envconfig:"SUBSCRIPTION_ID"`
}

// SQL defines the configuration for SQL database.
type SQL struct {
	Dialect     string `envconfig:"MQ_DB_DIALECT"`
	Host        string `envconfig:"MQ_DB_HOST"`
	Port        int    `envconfig:"MQ_DB_PORT"`
	Username    string `envconfig:"MQ_DB_USERNAME"`
	Password    string `envconfig:"MQ_DB_PASSWORD"`
	Db          string `envconfig:"MQ_DB_NAME"`
	QueryString string `envconfig:"DB_QUERYSTRING"`
}
