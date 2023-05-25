package app

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// database collection names
const (
	CollectionMints        = "mints"
	CollectionInvalidMints = "invalid_mints"
)

// Database is a wrapper around the mongo database
type Database struct {
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	// DB is the global database wrapper
	DB *Database
)

// Connect connects to the database
func (d *Database) Connect() error {
	log.Info("Connecting to database")
	client, err := mongo.Connect(d.ctx, options.Client().ApplyURI(Config.MongoDB.URI).SetWriteConcern(writeconcern.New(writeconcern.WMajority())))
	if err != nil {
		return err
	}
	d.db = client.Database(Config.MongoDB.Database)
	return nil
}

// Setup Indexes
func (d *Database) SetupIndexes() error {
	log.Info("Setting up indexes")

	// setup unique index for mints
	_, err := d.db.Collection(CollectionMints).Indexes().CreateOne(d.ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"transaction_hash": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err

	}

	// setup unique index for invalid mints
	_, err = d.db.Collection(CollectionInvalidMints).Indexes().CreateOne(d.ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"transaction_hash": 1},
		Options: options.Index().SetUnique(true),
	})

	// create multiple indexes in one call
	return nil
}

// Disconnect disconnects from the database
func (d *Database) Disconnect() error {
	log.Info("Disconnecting from database")
	err := d.db.Client().Disconnect(d.ctx)
	d.cancel()
	return err
}

// GetCollection gets a collection from the database
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

// NewDatabase creates a new database wrapper
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	DB = &Database{
		ctx:    ctx,
		cancel: cancel,
	}
	DB.Connect()
	DB.SetupIndexes()
}
