package member

import "go.mongodb.org/mongo-driver/mongo"

type GroupRepository interface {
}

type MongoGroupRepository struct {
	collection *mongo.Collection
}

func NewMongoGroupRepository(client *mongo.Client, dbName string, collectionName string) *MongoGroupRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoGroupRepository{
		collection: collection,
	}
}

