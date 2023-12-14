package services

import (
	"RedhaBena/indexer/internal/database"
	"context"
	"fmt"
)

var Client ServiceClient

type ServiceClient struct {
	neo4j *database.Neo4jClient
}

func InitServiceClient(ctx context.Context) error {
	Neo4jClient, err := database.NewNeo4jClient()
	if err != nil {
		return fmt.Errorf("failed to create Neo4j client: %v", err)
	}

	Client = *new(ServiceClient)
	Client.neo4j = Neo4jClient

	return nil
}
