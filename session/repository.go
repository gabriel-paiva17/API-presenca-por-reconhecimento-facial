package session

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	collection *mongo.Collection
}

func NewSessionRepository(client *mongo.Client, dbName string, collectionName string) *SessionRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &SessionRepository{
		collection: collection,
	}
}
