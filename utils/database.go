package utils

import (
	"fmt"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	fmt.Println("Conectado ao MongoDB!")
	return client, ctx, cancel, nil
}


