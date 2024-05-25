package user

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(client *mongo.Client, dbName string, collectionName string) *MongoUserRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoUserRepository{
		collection: collection,
	}
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}