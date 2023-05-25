package app

import (
	"context"
	"time"

	"github.com/dan13ram/wpokt-backend/models"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Database is a wrapper around the mongo database
type Database struct {
	db *mongo.Database
}

var (
	// DB is the global database wrapper
	DB *Database
)

// Connect connects to the database
func (d *Database) Connect(ctx context.Context) error {
	log.Info("Connecting to database")

	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(time.Duration(Config.MongoDB.TimeOutSecs)*time.Second))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Config.MongoDB.URI).SetWriteConcern(wcMajority))
	if err != nil {
		return err
	}
	d.db = client.Database(Config.MongoDB.Database)
	return nil
}

// Setup Indexes
func (d *Database) SetupIndexes() error {
	log.Info("Setting up indexes")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	// setup unique index for mints
	_, err := d.db.Collection(models.CollectionMints).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"transaction_hash": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err

	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	// setup unique index for invalid mints
	_, err = d.db.Collection(models.CollectionInvalidMints).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"transaction_hash": 1},
		Options: options.Index().SetUnique(true),
	})

	// create multiple indexes in one call
	return nil
}

// Disconnect disconnects from the database
func (d *Database) Disconnect() error {
	log.Info("Disconnecting from database")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	err := d.db.Client().Disconnect(ctx)
	return err
}

// GetCollection gets a collection from the database
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

// NewDatabase creates a new database wrapper
func InitDB(ctx context.Context) {
	DB = &Database{}
	DB.Connect(ctx)
	DB.SetupIndexes()
}
