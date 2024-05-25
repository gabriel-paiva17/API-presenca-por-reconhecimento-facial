package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"fmt"
	//"myproject/user"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func CreateServer() {

	// Conecta-se ao MongoDB
	client, cancel, err := connectDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Erro ao conectar ao MongoDB:", err)
	}
	defer func() {
		
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal("Erro ao desconectar do MongoDB:", err)
		}
		cancel() // Cancela o contexto para liberar recursos
	}()

	r := mux.NewRouter()

	s := &http.Server{
		Addr:         "192.168.0.9:8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	log.Fatal(s.ListenAndServe())

}


func connectDB(uri string) (*mongo.Client, context.CancelFunc, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	// Verifica se a conexão foi estabelecida com sucesso
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		cancel()
		return nil, nil, err
	}

	fmt.Println("Conectado ao MongoDB!")
	return client, cancel, nil
}