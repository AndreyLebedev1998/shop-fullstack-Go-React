package createkafkatopic

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func CreateKafkaTopic(broker string, topic string, partitions int, replication int) error {
	conn, err := kafka.Dial("tcp", broker)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer conn.Close()

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	})
}
