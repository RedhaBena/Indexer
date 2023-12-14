package database

import (
	"context"
	"fmt"

	"RedhaBena/indexer/internal/config"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jClient struct {
	driver *neo4j.DriverWithContext
	ctx    context.Context
}

func NewNeo4jClient() (*Neo4jClient, error) {
	uri := fmt.Sprintf("neo4j://%s", config.GlobalConfig.DatabaseConfig.Host)

	user := config.GlobalConfig.DatabaseConfig.User
	password := config.GlobalConfig.DatabaseConfig.Pass

	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(user, password, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Neo4j: %v", err)
	}

	ctx := context.Background()

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Neo4j connectivity: %v", err)
	}

	client := new(Neo4jClient)
	client.driver = &driver
	client.ctx = ctx

	if err := client.createIndexes(); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	return client, nil
}

func (client *Neo4jClient) createIndexes() error {
	queryArticle := "CREATE CONSTRAINT IF NOT EXISTS FOR (a:Article) REQUIRE a.id IS UNIQUE;"
	_, err := neo4j.ExecuteQuery(client.ctx, *client.driver, queryArticle, nil,
		neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return err
	}

	queryAuthor := "CREATE CONSTRAINT IF NOT EXISTS FOR (a:Author) REQUIRE a.id IS UNIQUE;"
	_, err = neo4j.ExecuteQuery(client.ctx, *client.driver, queryAuthor, nil,
		neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))

	return err
}
