// main.go
package services

import (
	"RedhaBena/indexer/internal/config"
	"RedhaBena/indexer/internal/logger"
	"RedhaBena/indexer/internal/models"
	"RedhaBena/indexer/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

func (client *ServiceClient) ReadFile(ctx context.Context) {
	file, err := os.Open(config.GlobalConfig.FileConfig.LocalPath)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	var articles []models.Article
	number_articles := 0

	// Create a buffered reader
	reader := utils.NewCustomReader(file)

	// Create a JSON decoder with a custom Unmarshaler for numbers
	decoder := json.NewDecoder(reader)

	// read open bracket
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	// Record the start time
	startTime := time.Now()
	logger.GlobalLogger.Debug("Starting to decode")

	var waitArticles sync.WaitGroup
	articleChannel := make(chan []models.Article, 100) // Max of 100 Batch in the queue

	// Goroutine to process articles
	go func() {
		defer close(articleChannel)
		for batch := range articleChannel {
			waitArticles.Add(1)
			defer waitArticles.Done()
			err := client.neo4j.InsertArticles(batch)
			if err != nil {
				logger.GlobalLogger.Error("Failed to create articles", zap.Error(err))
				return
			}
		}
	}()

	for decoder.More() {
		select {
		case <-ctx.Done():
			logger.GlobalLogger.Debug("Stopping decoder loop - context done")
			return
		default:
			// Parse the line into the Data struct
			var article models.Article

			err := decoder.Decode(&article)
			if err != nil {
				log.Fatal(err)
			}

			for i, author := range article.Authors {
				if author.Id == nil {
					id := fmt.Sprintf("%s%x", article.Id, i)
					article.Authors[i].Id = &id
				}
			}
			articles = append(articles, article)

			if uint(len(articles)) >= config.GlobalConfig.IndexerConfig.BatchSize {
				number_articles = number_articles + len(articles)
				logger.GlobalLogger.Debug("Inserting articles", zap.Int("size", len(articles)), zap.Int("total", number_articles))
				articleChannel <- articles
				articles = []models.Article{}
			}
		}
	}

	// Process the remaining articles
	if len(articles) > 0 {
		articleChannel <- articles
	}
	number_articles = number_articles + len(articles)
	logger.GlobalLogger.Debug("Inserting articles", zap.Int("size", len(articles)), zap.Int("total", number_articles))
	waitArticles.Wait()

	sinceTime := time.Since(startTime)
	logger.GlobalLogger.Info("End of insert", zap.Duration("duration", sinceTime))
}
