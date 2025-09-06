package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"

	"profile-golang/common/cache"
	"profile-golang/common/models"
	"profile-golang/common/redis"
	"profile-golang/config"
	"profile-golang/reader/handler"
)

type Server struct {
	config      *config.Config
	port        int
	redis       *redis.Client
	cache       *cache.SimpleCache
	router      *gin.Engine
	kafkaReader *kafka.Reader
}

func NewServer(cfg *config.Config, port int) *Server {
	redisClient, err := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	simpleCache := cache.NewSimpleCache(cfg.CacheCapacity)

	router := gin.Default()

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		GroupID:  "cache-readers",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	server := &Server{
		config:      cfg,
		port:        port,
		redis:       redisClient,
		cache:       simpleCache,
		router:      router,
		kafkaReader: kafkaReader,
	}

	handler.RegisterRoutes(router, server.cache, redisClient)

	return server
}

func (s *Server) Start() error {
	go s.consumeKafka()

	go s.startCleanupJob()

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Reader server starting on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) consumeKafka() {
	log.Println("Starting Kafka consumer...")

	for {
		msg, err := s.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading Kafka message: %v", err)
			continue
		}

		var kafkaMsg models.KafkaMessage
		if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
			log.Printf("Error parsing Kafka message: %v", err)
			continue
		}

		switch kafkaMsg.Type {
		case "update":
			s.cache.Set(kafkaMsg.Key, kafkaMsg.Value)
			log.Printf("Updated cache from Kafka: %s = %s", kafkaMsg.Key, kafkaMsg.Value)
		case "delete":
			s.cache.Delete(kafkaMsg.Key)
			log.Printf("Deleted cache from Kafka: %s", kafkaMsg.Key)
		}
	}
}

func (s *Server) startCleanupJob() {
	c := cron.New()

	// Daily cache cleanup at midnight
	_, err := c.AddFunc(s.config.CacheCleanupCron, func() {
		log.Println("Starting cache cleanup...")
		s.cache.Cleanup(24 * time.Hour)
		log.Println("Cache cleanup completed")
	})

	if err != nil {
		log.Printf("Error scheduling cleanup job: %v", err)
		return
	}

	c.Start()
}
