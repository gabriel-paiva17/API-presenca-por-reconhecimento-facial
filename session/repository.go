package session

import (
	"context"
	"time"

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

// EndSession atualiza o campo EndedAt de uma sessão específica
func (r *SessionRepository) EndSession(ctx context.Context, session *Session) error {
	filter := bson.M{
		"_id": session.ID,
	}


	update := bson.M{
		"$set": bson.M{
			"endedAt": time.Now().Format(time.RFC3339),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GET /grupos/{nome-do-grupo}/sessoes/finalizadas
// e 
// GET /grupos/{nome-do-grupo}/sessoes/em-andamento

func (r *SessionRepository) findAllWithFilter(ctx context.Context, filter bson.M) (*GetManySessionsResponse, error) {
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var sessions []SessionByName
    for cursor.Next(ctx) {
        var session Session
        if err := cursor.Decode(&session); err != nil {
            return nil, err
        }
        sessions = append(sessions, SessionByName{Name: session.Name})
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return &GetManySessionsResponse{Sessions: sessions}, nil
}

func (r *SessionRepository) FindAllActiveSessions(ctx context.Context, groupName string, userID string) (*GetManySessionsResponse, error) {
    filter := bson.M{
        "groupName": groupName,
        "createdBy": userID,
        "endedAt":   bson.M{"$eq": ""},
    }
    return r.findAllWithFilter(ctx, filter)
}

func (r *SessionRepository) FindAllEndedSessions(ctx context.Context, groupName string, userID string) (*GetManySessionsResponse, error) {
    filter := bson.M{
        "groupName": groupName,
        "createdBy": userID,
        "endedAt":   bson.M{"$ne": ""},
    }
    return r.findAllWithFilter(ctx, filter)
}

// DELETE /grupos/{nome-do-grupo}/sessoes/{nome-da-sessao}

func (r *SessionRepository) DeleteOneSession(ctx context.Context, groupName, createdBy, sessionName string) error {
    filter := bson.M{
        "group_name":  groupName,
        "created_by":  createdBy,
        "session_name": sessionName,
    }
    _, err := r.collection.DeleteOne(ctx, filter)
    return err
}