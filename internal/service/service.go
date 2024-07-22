package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/kuromii5/messagio/internal/models"
	le "github.com/kuromii5/messagio/pkg/logger/l_err"
)

type MessageService struct {
	logger *slog.Logger

	messageSaver MessageSaver
	statsLoader  StatsLoader

	producer   sarama.SyncProducer
	consumer   sarama.Consumer
	kafkaTopic string
}

type MessageSaver interface {
	Save(ctx context.Context, message string) (int32, error)
	MarkAsProcessed(ctx context.Context, messageID int32) error
}
type StatsLoader interface {
	LoadStats(ctx context.Context) (int32, error)
}

func NewService(
	logger *slog.Logger,
	messageSaver MessageSaver,
	statsLoader StatsLoader,
	producer sarama.SyncProducer,
	consumer sarama.Consumer,
	kafkaTopic string,
) *MessageService {
	return &MessageService{
		logger:       logger,
		messageSaver: messageSaver,
		statsLoader:  statsLoader,
		producer:     producer,
		consumer:     consumer,
		kafkaTopic:   kafkaTopic,
	}
}

func (m *MessageService) SendMessage(ctx context.Context, msg string) (models.SendMessageResponse, error) {
	m.logger.Info("sending message...")

	// Save the message to the database
	messageID, err := m.messageSaver.Save(ctx, msg)
	if err != nil {
		m.logger.Error("failed to save message in postgres", le.Err(err))

		return models.SendMessageResponse{
			StatusCode: 500,
			Message:    "Failed to save message",
		}, err
	}

	// Send the message to Kafka
	kafkaMessage := &sarama.ProducerMessage{
		Topic: m.kafkaTopic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err = m.producer.SendMessage(kafkaMessage)
	if err != nil {
		m.logger.Error("failed to send message to Kafka", le.Err(err))

		return models.SendMessageResponse{
			StatusCode: 500,
			Message:    "Failed to send message to Kafka",
		}, err
	}

	// Update the message as processed in the database
	err = m.messageSaver.MarkAsProcessed(ctx, messageID)
	if err != nil {
		m.logger.Error("message sent to Kafka, but failed to update status in database", le.Err(err))

		return models.SendMessageResponse{
			StatusCode: 500,
			Message:    "Message sent to Kafka, but failed to update status in database",
		}, err
	}

	m.logger.Info("Successfully sent message", slog.Int("message_id", int(messageID)))

	// Respond
	return models.SendMessageResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Message sent successfully. ID: %d", messageID),
	}, nil
}

func (m *MessageService) GetMessages(ctx context.Context) (models.GetMessagesResponse, error) {
	m.logger.Info("fetching messages...")

	partitionConsumer, err := m.consumer.ConsumePartition(m.kafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		m.logger.Error("Failed to consume topic", slog.String("topic", m.kafkaTopic), le.Err(err))

		return models.GetMessagesResponse{}, err
	}
	defer partitionConsumer.Close()

	var messages []string
	msgChan := partitionConsumer.Messages()

	// Consume messages for a limited period or number of messages
	timeout := time.After(5 * time.Second)
	for {
		select {
		case msg := <-msgChan:
			m.logger.Info("Message received", slog.String("message", string(msg.Value)))

			messages = append(messages, string(msg.Value))

		case <-timeout:
			m.logger.Info("Successfully fetched messages from topic", slog.String("topic", m.kafkaTopic))

			return models.GetMessagesResponse{Messages: messages}, nil

		case err := <-partitionConsumer.Errors():
			m.logger.Error("Error consuming messages", le.Err(err))

			return models.GetMessagesResponse{}, err
		}
	}
}

func (m *MessageService) GetStats(ctx context.Context) (models.GetStatsResponse, error) {
	m.logger.Info("fetching statistics...")

	processedCount, err := m.statsLoader.LoadStats(ctx)
	if err != nil {
		m.logger.Error("failed to fetch statistics", le.Err(err))

		return models.GetStatsResponse{}, err
	}

	m.logger.Info("successfully fetched statistics", slog.Int("processed_messages", int(processedCount)))

	return models.GetStatsResponse{ProcessedCount: processedCount}, nil
}
