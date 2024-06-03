package group

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepository struct {
	collection *mongo.Collection
}

func NewGroupRepository(client *mongo.Client, dbName string, collectionName string) *GroupRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &GroupRepository{
		collection: collection,
	}
}

func (r *GroupRepository) FindAllGroupsByUserID(userID string, ctx context.Context) ([]Group, error) {
	
	var groups []Group

	filter := bson.M{"createdBy": userID}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GroupRepository) FindOneByNameAndCreator(ctx context.Context, name string, createdBy string) (*Group, bool) {

	filter := bson.M{"name": name, "createdBy": createdBy}
    existingGroup := &Group{}
	
	err := r.collection.FindOne(ctx, filter).Decode(existingGroup)

	if err != nil {

		return nil, false

	}

	return existingGroup, true 

}

// POST /group

func (r *GroupRepository) CreateGroup(ctx context.Context, group *Group) error {

	_, found :=  r.FindOneByNameAndCreator(ctx, group.Name, group.CreatedBy)
	
	if found { 
        return ErrNameAlreadyExists
    }

    _, err := r.collection.InsertOne(ctx, group)
    
	return err

}