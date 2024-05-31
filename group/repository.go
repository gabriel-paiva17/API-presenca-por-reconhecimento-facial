package group

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
)

type GroupRepository interface {

	CreateUser(ctx context.Context, group *Group) error

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

// func (r *MongoGroupRepository)

func (r *MongoGroupRepository) CreateUser(ctx context.Context, group *Group) error {






	return nil



}