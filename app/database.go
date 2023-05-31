package app

import (
	"context"
	"time"

	"github.com/dan13ram/wpokt-backend/models"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
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
	log.Debug("Connecting to database")

	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(time.Duration(Config.MongoDB.TimeOutSecs)*time.Second))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Config.MongoDB.URI).SetWriteConcern(wcMajority))
	if err != nil {
		return err
	}
	d.db = client.Database(Config.MongoDB.Database)

	log.Debug("Connected to database")
	return nil
}

// Setup Indexes
func (d *Database) SetupIndexes() error {
	log.Debug("Setting up indexes")

	// setup unique index for mints
	log.Debug("Setting up indexes for mints")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	_, err := d.db.Collection(models.CollectionMints).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil {
		return err
	}

	// setup unique index for invalid mints
	log.Debug("Setting up indexes for invalid mints")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	_, err = d.db.Collection(models.CollectionInvalidMints).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil {
		return err
	}

	// setup unique index for burns
	log.Debug("Setting up indexes for burns")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	_, err = d.db.Collection(models.CollectionBurns).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}, {Key: "log_index", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil {
		return err
	}

	log.Debug("Set up indexes")

	return nil
}

// Disconnect disconnects from the database
func (d *Database) Disconnect() error {
	log.Debug("Disconnecting from database")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	err := d.db.Client().Disconnect(ctx)
	log.Debug("Disconnected from database")
	return err
}

// GetCollection gets a collection from the database
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

// NewDatabase creates a new database wrapper
func InitDB(ctx context.Context) {
	DB = &Database{}
	err := DB.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = DB.SetupIndexes()
	if err != nil {
		log.Fatal(err)
	}
}
