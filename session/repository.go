package session

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

func (r *SessionRepository) FindOneSession(ctx context.Context, groupName string, userID string, sessionName string) (*Session, bool) {
	
	filter := bson.M{
		"groupName": groupName,
		"createdBy": userID,
		"name":      sessionName,
	}

	existingSession := &Session{}

	err := r.collection.FindOne(ctx, filter).Decode(existingSession)
	if err != nil {

		return nil, false

	}

	// Sessão encontrada
	return existingSession, true
}

// POST /grupos/{nome-do-grupo}/sessoes/iniciar

func (r *SessionRepository) StartNewSession(ctx context.Context, session *Session) error {

	_, found := r.FindOneSession(ctx, session.GroupName, session.CreatedBy, session.Name)

	if found {

		return ErrSessionAlreadyExists

	}

	_, err := r.collection.InsertOne(ctx, session)

	return err

}

func (r *SessionRepository) UpdateMembers(ctx context.Context, session *Session, newMembers []SessionMember) error {

	// Cria o filtro para encontrar a sessão
	filter := bson.M{
		"_id": session.ID,
	}

	// Define a atualização usando a operação $set
	update := bson.M{
		"$set": bson.M{
			"members": newMembers,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}
