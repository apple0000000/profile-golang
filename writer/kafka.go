package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"profile-golang/common/models"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}

	return &Producer{
		writer: w,
		topic:  topic,
	}, nil
}

func (p *Producer) SendMessage(message models.KafkaMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(message.Key),
			Value: jsonData,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Sent Kafka message: %s %s", message.Type, message.Key)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
