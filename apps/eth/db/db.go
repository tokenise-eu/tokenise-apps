package db

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// DB holds a database instance and it's associated client.
type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Connect will connect to the database and return an instance of it.
func Connect(user string, pass string) (*DB, error) {
	client, err := mongo.NewClient("mongodb://" + user + ":" + pass + "@localhost:27017")
	if err != nil {
		return nil, err
	}

	if err := client.Connect(context.Background()); err != nil {
		return nil, err
	}

	return &DB{
		Client: client,
		DB:     client.Database("contract"),
	}, nil
}

// AddUser is mostly for posterity right now
// This should ideally be done from Onfido and not from the contract.
func (db *DB) AddUser(from []byte) error {
	// Unmarshal first, but this structure will be apparent later
	doc := bson.D{{"from", from}}
	res, err := db.DB.Collection("users").InsertOne(context.Background(), doc)
	fmt.Println(res)
	return err
}

// AddTx will add a transaction to the transactions database
func (db *DB) AddTx(from string, to string, value string) error {
	doc := bson.D{{"from", from}, {"to", to}, {"value", value}}
	res, err := db.DB.Collection("txs").InsertOne(context.Background(), doc)
	fmt.Println(res)
	return err
}

// Close the connection to the database.
func (db *DB) Close() error {
	return db.Client.Disconnect(context.Background())
}
