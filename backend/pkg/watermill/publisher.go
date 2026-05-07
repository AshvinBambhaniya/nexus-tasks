// Package watermill provides message queue functionality using Watermill.
package watermill

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/cli/workers"
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/database"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill-googlecloud/v2/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

// IPublisher is an interface for publishing messages to a topic.
type IPublisher interface {
	Publish(topic string, handle workers.Handler) error
}

// Publisher is an implementation of Publisher interface using Watermill.
type Publisher struct {
	publisher message.Publisher
}

// InitPublisher initializes a new Publisher.
func InitPublisher(cfg config.AppConfig, isDeadLetterQ bool) (*Publisher, error) {
	logger = watermill.NewStdLogger(cfg.MQ.Debug, cfg.MQ.Track)
	if isDeadLetterQ {
		return initSQLPub(cfg)
	}
	switch cfg.MQ.Dialect {
	case "amqp":
		return initAmqpPub(cfg)

	case "redis":
		return initRedisPub(cfg)

	case "kafka":
		return initKafkaPub(cfg)

	case "googlecloud":
		return initGoogleCloudPub(cfg)

	case "sql":
		return initSQLPub(cfg)
	default:
		return &Publisher{}, nil
	}
}

// Publish sends a message into the queue using a topic name.
func (wp *Publisher) Publish(topic string, handle workers.Handler) error {
	// if broker is not set then return nil
	if wp.publisher == nil {
		return nil
	}
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)

	err := enc.Encode(&handle)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), network.Bytes())
	err = wp.publisher.Publish(topic, msg)
	return err
}

func initAmqpPub(cfg config.AppConfig) (*Publisher, error) {
	amqpConfig := amqp.NewDurableQueueConfig(cfg.MQ.Amqp.AmqpURL)
	publisher, err := amqp.NewPublisher(amqpConfig, logger)
	return &Publisher{publisher: publisher}, err
}

func initRedisPub(cfg config.AppConfig) (*Publisher, error) {
	pubClient := redis.NewClient(&redis.Options{
		Addr:     cfg.MQ.Redis.RedisURL,
		Username: cfg.MQ.Redis.UserName,
		Password: cfg.MQ.Redis.Password,
	})
	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client:     pubClient,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	return &Publisher{publisher: publisher}, err
}

func initKafkaPub(cfg config.AppConfig) (*Publisher, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               cfg.MQ.Kafka.KafkaBroker,
			Marshaler:             kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: kafka.DefaultSaramaSyncPublisherConfig(),
		},
		logger,
	)
	return &Publisher{publisher: publisher}, err
}

func initGoogleCloudPub(cfg config.AppConfig) (*Publisher, error) {
	publisher, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID:      cfg.MQ.GoogleCloud.ProjectID,
		ConnectTimeout: 10 * time.Second,
		PublishTimeout: 10 * time.Second,
		Marshaler:      googlecloud.DefaultMarshalerUnmarshaler{},
	}, logger)

	return &Publisher{publisher: publisher}, err
}

func initSQLPub(cfg config.AppConfig) (*Publisher, error) {
	switch cfg.MQ.SQL.Dialect {
	case "postgres":
		return initPostgresPub(cfg)
	case "mysql":
		return initMysqlPub(cfg)
	default:
		return nil, nil
	}
}

func initPostgresPub(cfg config.AppConfig) (*Publisher, error) {
	db, err := database.PostgresDBConnection(cfg.MQ.SQL)
	if err != nil {
		return nil, err
	}
	publisher, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			SchemaAdapter:        database.PostgreSQLSchema{},
			AutoInitializeSchema: true,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}
	return &Publisher{publisher: publisher}, nil
}

func initMysqlPub(cfg config.AppConfig) (*Publisher, error) {
	db, err := database.MysqlDBConnection(cfg.MQ.SQL)
	if err != nil {
		return nil, err
	}
	publisher, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			SchemaAdapter:        database.MySQLSchema{},
			AutoInitializeSchema: true,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{publisher: publisher}, nil
}
