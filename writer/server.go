package writer

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"profile-golang/common/redis"
	"profile-golang/config"
	"profile-golang/writer/handler"
	"profile-golang/writer/kafka"
)

type Server struct {
	config    *config.Config
	port      int
	redis     *redis.Client
	kafkaProd *kafka.Producer
	router    *gin.Engine
}

func NewServer(cfg *config.Config, port int) *Server {
	redisClient, err := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	router := gin.Default()

	server := &Server{
		config:    cfg,
		port:      port,
		redis:     redisClient,
		kafkaProd: kafkaProducer,
		router:    router,
	}

	handler.RegisterRoutes(router, redisClient, kafkaProducer)

	return server
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Writer server starting on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) Close() {
	s.redis.Close()
	s.kafkaProd.Close()
}
