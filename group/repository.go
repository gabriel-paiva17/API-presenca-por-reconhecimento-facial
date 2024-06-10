package group

import (
	"context"
	"myproject/cv"

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

// GET /grupos

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

func (r *GroupRepository) FindOneByNameAndCreator(ctx context.Context, groupName string, createdBy string) (*Group, bool) {

	filter := bson.M{"name": groupName, "createdBy": createdBy}
	existingGroup := &Group{}

	err := r.collection.FindOne(ctx, filter).Decode(existingGroup)

	if err != nil {

		return nil, false

	}

	return existingGroup, true

}

// POST /grupos/criar

func (r *GroupRepository) CreateGroup(ctx context.Context, group *Group) error {

	_, found := r.FindOneByNameAndCreator(ctx, group.Name, group.CreatedBy)

	if found {
		return ErrNameAlreadyExists
	}

	_, err := r.collection.InsertOne(ctx, group)

	return err

}

// POST /grupos/{nome-do-grupo}/detalhes/adicionar

func (r *GroupRepository) AddMemberToGroup(ctx context.Context, groupName, createdBy string, newMember *Member) (*Member, error) {
	// Verificar se já existe um membro com o mesmo nome no grupo
	filter := bson.M{
		"name":      groupName,
		"createdBy": createdBy,
		"members": bson.M{
			"$elemMatch": bson.M{"name": newMember.Name},
		},
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrNameAlreadyExists
	}

	// Obter o grupo
	group, found := r.FindOneByNameAndCreator(ctx, groupName, createdBy)
	if !found {
		return nil, ErrGroupNotFound
	}

	// Verificar se a face do novo membro é a mesma de algum membro existente no grupo
	for _, existingMember := range group.Members {
		samePerson, err := cv.CompareFaces(newMember.Face, existingMember.Face)
		if err != nil {
			return nil, err
		}
		if samePerson {
			return nil, ErrFaceAlreadyUsed
		}
	}

	// Adicionar o novo membro ao array de membros do grupo
	update := bson.M{"$push": bson.M{"members": newMember}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"name": groupName, "createdBy": createdBy}, update)
	if err != nil {
		return nil, err
	}
	return newMember, nil
}