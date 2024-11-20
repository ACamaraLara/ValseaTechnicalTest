package mongodb

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct that stores a connection to MongoDataBase.
type MongoDBClient struct {
	Client      *mongo.Client
	Config      *MongoConfig
	Collections map[string]*mongo.Collection
}

func NewMongoDBClient(cfg *MongoConfig) *MongoDBClient {
	return &MongoDBClient{Config: cfg}
}

// ConnectMongoClient connects the initialized object to the given database.
func (mb *MongoDBClient) ConnectMongoClient(ctx context.Context) error {

	clientOptions := options.Client().ApplyURI(mb.Config.GetURL())

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB url: %v", err)
	}

	mb.Client = client

	if !mb.checkConnection(ctx) {
		return fmt.Errorf("mongodb connection check failed")
	}

	log.Info().Msg("Connected to MongoDB on url " + mb.Config.GetURL())

	return nil
}

func (mb *MongoDBClient) GetCollections(colls []string) error {
	mb.Collections = make(map[string]*mongo.Collection)
	for _, coll := range colls {
		mb.Collections[coll] = mb.Client.Database(mb.Config.DbName).Collection(coll)
	}
	return nil
}

// Function to check if the client is successfully connected to database.
func (mb *MongoDBClient) checkConnection(ctx context.Context) bool {

	if mb.Client == nil {
		return false
	}

	err := mb.Client.Ping(ctx, nil)
	if err != nil {
		log.Error().Msg("Error pinging mongoDB:" + err.Error())
	}

	return err == nil
}
