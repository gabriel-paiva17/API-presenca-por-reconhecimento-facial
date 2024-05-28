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

func (r *MongoUserRepository) FindOneByEmail (ctx context.Context, email string) (*User, bool) {

	filter := bson.M{"email": email}
    existingUser := &User{}
	
	// caso nao tenha erro, significa que a busca funcionou e existe
	// um user cadastrado com esse email
    err := r.collection.FindOne(ctx, filter).Decode(existingUser)

	if err != nil {

		return nil, false 

	}

	return existingUser, true

}

var ErrEmailAlreadyExists = errors.New("email already used")
 

/////////////////////////
// POST /auth/register //
/////////////////////////


func (r *MongoUserRepository) CreateUser(ctx context.Context, user *User) error {
	
	_, exists :=  r.FindOneByEmail(ctx, user.Email)
	
	if exists { 
        return ErrEmailAlreadyExists
    }

    _, err := r.collection.InsertOne(ctx, user)
    
	return err

}