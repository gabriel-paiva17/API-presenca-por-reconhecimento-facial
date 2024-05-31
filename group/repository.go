package group

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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


func (r *MongoGroupRepository) FindOneByName(ctx context.Context, name string) (*Group, bool) {

	filter := bson.M{"name": name}
    existingGroup := &Group{}
	
	err := r.collection.FindOne(ctx, filter).Decode(existingGroup)

	if err != nil {

		return nil, false

	}

	return existingGroup, true 

}

// POST /group

func (r *MongoGroupRepository) CreateGroup(ctx context.Context, group *Group) error {

	_, found :=  r.FindOneByName(ctx, group.Name)
	
	if found { 
        return ErrNameAlreadyExists
    }

    _, err := r.collection.InsertOne(ctx, group)
    
	return err

}