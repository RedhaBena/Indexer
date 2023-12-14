package database

import (
	"RedhaBena/indexer/internal/logger"
	"RedhaBena/indexer/internal/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

func (client *Neo4jClient) InsertArticles(articles []models.Article) error {
	logger.GlobalLogger.Debug("Inserting articles", zap.Int("size", len(articles)))

	query := `
		UNWIND $articles AS article
		MERGE (a:Article {id: article.id})
		SET a.title = article.title
		WITH a, article.authors AS authors, article.references AS references
		UNWIND authors AS authorData
		MERGE (author:Author {id: authorData.id})
		ON CREATE SET author.name = authorData.name
		MERGE (author)-[:AUTHORED]->(a)
		WITH a, references
		UNWIND references AS reference
		MERGE (b:Article {id: reference})
		MERGE (a)-[:CITE]->(b)
	`
	data := map[string]interface{}{
		"articles": func() []map[string]interface{} {
			result := make([]map[string]interface{}, len(articles))
			for i, article := range articles {
				result[i] = map[string]interface{}{
					"id":         article.Id,
					"title":      article.Title,
					"references": article.References,
					"authors": func() []map[string]interface{} {
						authorData := make([]map[string]interface{}, len(article.Authors))
						for j, author := range article.Authors {
							authorData[j] = map[string]interface{}{
								"id":   *author.Id,
								"name": author.Name,
							}
						}
						return authorData
					}()}
			}
			return result
		}(),
	}

	_, err := neo4j.ExecuteQuery(client.ctx, *client.driver, query, data,
		neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

	return err
}
