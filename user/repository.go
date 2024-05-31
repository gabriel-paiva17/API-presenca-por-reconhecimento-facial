package user

import (
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(client *mongo.Client, dbName string, collectionName string) *UserRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) FindOneByEmail (ctx context.Context, email string) (*User, bool) {

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

/////////////////////////
// POST /auth/register //
/////////////////////////


func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	
	_, found :=  r.FindOneByEmail(ctx, user.Email)
	
	if found { 
        return ErrEmailAlreadyExists
    }

    _, err := r.collection.InsertOne(ctx, user)
    
	return err

}