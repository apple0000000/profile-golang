package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"profile-golang/common/cache"
	"profile-golang/common/models"
	"profile-golang/common/redis"
)

func RegisterRoutes(router *gin.Engine, cache *cache.SimpleCache, redisClient *redis.Client) {
	router.GET("/read", func(c *gin.Context) {
		keysParam := c.Query("keys")
		if keysParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "keys parameter is required"})
			return
		}

		keys := strings.Split(keysParam, ",")
		if len(keys) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 100 keys allowed per request"})
			return
		}

		response := readItems(keys, cache, redisClient)
		c.JSON(http.StatusOK, response)
	})
}

func readItems(keys []string, cache *cache.SimpleCache, redisClient *redis.Client) models.ReadResponse {
	var response models.ReadResponse
	response.Items = make([]models.ReadResponseItem, len(keys))

	// First check memory cache
	for i, key := range keys {
		if value, found := cache.Get(key); found {
			response.Items[i] = models.ReadResponseItem{
				Key:    key,
				Value:  value,
				Source: "memory",
				Exists: true,
			}
		}
	}

	// Find keys not in memory cache
	var missingKeys []string
	var missingIndices []int
	for i, item := range response.Items {
		if item.Source == "" {
			missingKeys = append(missingKeys, keys[i])
			missingIndices = append(missingIndices, i)
		}
	}

	if len(missingKeys) > 0 {
		// Get from Redis in batch
		redisValues, err := redisClient.MGet(missingKeys)
		if err != nil {
			log.Printf("Error reading from Redis: %v", err)
			for _, idx := range missingIndices {
				response.Items[idx] = models.ReadResponseItem{
					Key:    keys[idx],
					Value:  "",
					Source: "error",
					Exists: false,
				}
			}
			return response
		}

		// Process Redis results
		for j, idx := range missingIndices {
			redisVal := redisValues[j]
			if redisVal != nil {
				valueStr := redisVal.(string)
				response.Items[idx] = models.ReadResponseItem{
					Key:    missingKeys[j],
					Value:  valueStr,
					Source: "redis",
					Exists: true,
				}
				// Update memory cache
				cache.Set(missingKeys[j], valueStr)
			} else {
				response.Items[idx] = models.ReadResponseItem{
					Key:    missingKeys[j],
					Value:  "",
					Source: "not_found",
					Exists: false,
				}
			}
		}
	}

	return response
}
