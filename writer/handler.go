package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"profile-golang/common/models"
	"profile-golang/common/redis"
	"profile-golang/writer/kafka"
)

func RegisterRoutes(router *gin.Engine, redisClient *redis.Client, kafkaProd *kafka.Producer) {
	router.POST("/write", func(c *gin.Context) {
		var request models.WriteRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(request.Items) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "items cannot be empty"})
			return
		}

		if len(request.Items) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 100 items allowed per request"})
			return
		}

		if err := writeItems(request.Items, redisClient, kafkaProd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}

func writeItems(items []models.Item, redisClient *redis.Client, kafkaProd *kafka.Producer) error {
	if err := redisClient.MSet(items); err != nil {
		return err
	}

	// Send Kafka messages for cache synchronization
	for _, item := range items {
		message := models.KafkaMessage{
			Type:  "update",
			Key:   item.Key,
			Value: item.Value,
		}
		if err := kafkaProd.SendMessage(message); err != nil {
			log.Printf("Failed to send Kafka message for key %s: %v", item.Key, err)
		}
	}

	return nil
}
