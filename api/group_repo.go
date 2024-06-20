package api

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

// DELETE /grupos/{nome-do-grupo}/deletar

func (r *GroupRepository) DeleteOneGroup(ctx context.Context, groupName string, createdBy string) error {

    filter := bson.M{"name": groupName, "createdBy": createdBy}

    _, err := r.collection.DeleteOne(ctx, filter)
    return err
}

// DELETE /grupos/deletar

func (r *GroupRepository) DeleteAllGroupsFromUser(ctx context.Context, createdBy string) error {

    filter := bson.M{"createdBy": createdBy}

    _, err := r.collection.DeleteMany(ctx, filter)
    return err
}

// DELETE /grupos/{nome-do-grupo}/detalhes/{nome-do-membro}/deletar

func (r *GroupRepository) RemoveOneMemberFromGroup(ctx context.Context, groupName, userID, memberName string) error {
    group, found := r.FindOneByNameAndCreator(ctx, groupName, userID)
    if !found {
        return ErrGroupNotFound
    }
    
	for i := range group.Members {
       
		if group.Members[i].Name == memberName {
            
			group.Members = append(group.Members[:i], group.Members[i+1:]...)
			
            return r.UpdateMembers(ctx, group, group.Members)
      
		}
    }
   
	return ErrMemberNotFound
}

// DELETE /grupos/{nome-do-grupo}/detalhes/deletar-membros

func (r *GroupRepository) RemoveAllMembersFromGroup(ctx context.Context, groupName, userID string) error {
    group, found := r.FindOneByNameAndCreator(ctx, groupName, userID)
    if !found {
        return ErrGroupNotFound
    }
    
	return r.UpdateMembers(ctx, group, []Member{})
}

func (r *GroupRepository) UpdateMembers(ctx context.Context, group *Group, newMembers []Member) error {

	filter := bson.M{
		"_id": group.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"members": newMembers,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}