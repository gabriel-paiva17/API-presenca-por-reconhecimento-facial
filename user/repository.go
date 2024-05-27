package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
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

var ErrEmailAlreadyExists = errors.New("email already used")
 
func (r *MongoUserRepository) CreateUser(ctx context.Context, user *User) error {
	
	filter := bson.M{"email": user.Email}
    existingUser := &User{}

    err := r.collection.FindOne(ctx, filter).Decode(existingUser)
    
	// caso nao tenha erro, significa que a busca funcionou e existe
	// um user cadastrado com esse email
	if err == nil { 
        return ErrEmailAlreadyExists
    }

    _, err = r.collection.InsertOne(ctx, user)
    return err
}