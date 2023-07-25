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

type Database interface {
	Connect(ctx context.Context) error
	SetupIndexes() error
	Disconnect() error
	InsertOne(collection string, data interface{}) error
	FindOne(collection string, filter interface{}, result interface{}) error
	FindMany(collection string, filter interface{}, result interface{}) error
	UpdateOne(collection string, filter interface{}, update interface{}) error
	UpsertOne(collection string, filter interface{}, update interface{}) error
}

// mongoDatabase is a wrapper around the mongo database
type mongoDatabase struct {
	db *mongo.Database
}

var (
	DB Database
)

// Connect connects to the database
func (d *mongoDatabase) Connect(ctx context.Context) error {
	log.Debug("[DB] Connecting to database")
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(time.Duration(Config.MongoDB.TimeOutSecs)*time.Second))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Config.MongoDB.URI).SetWriteConcern(wcMajority))
	if err != nil {
		return err
	}
	d.db = client.Database(Config.MongoDB.Database)

	log.Info("[DB] Connected to mongo database: ", Config.MongoDB.Database)
	return nil
}

// Setup Indexes
func (d *mongoDatabase) SetupIndexes() error {
	log.Debug("[DB] Setting up indexes")

	// setup unique index for mints
	log.Debug("[DB] Setting up indexes for mints")
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
	log.Debug("[DB] Setting up indexes for invalid mints")
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
	log.Debug("[DB] Setting up indexes for burns")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	_, err = d.db.Collection(models.CollectionBurns).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}, {Key: "log_index", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil {
		return err
	}

	log.Info("[DB] Indexes setup")

	return nil
}

// Disconnect disconnects from the database
func (d *mongoDatabase) Disconnect() error {
	log.Debug("[DB] Disconnecting from database")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	err := d.db.Client().Disconnect(ctx)
	log.Info("[DB] Disconnected from database")
	return err
}

// InitDB creates a new database wrapper
func InitDB(ctx context.Context) {
	DB = &mongoDatabase{}
	err := DB.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = DB.SetupIndexes()
	if err != nil {
		log.Fatal(err)
	}
}

// method for insert single value in a collection
func (d *mongoDatabase) InsertOne(collection string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	_, err := d.db.Collection(collection).InsertOne(ctx, data)
	return err
}

// method for find single value in a collection
func (d *mongoDatabase) FindOne(collection string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	err := d.db.Collection(collection).FindOne(ctx, filter).Decode(result)
	return err
}

// method for find multiple values in a collection
func (d *mongoDatabase) FindMany(collection string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	cursor, err := d.db.Collection(collection).Find(ctx, filter)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, result)
	return err
}

//method for update single value in a collection
func (d *mongoDatabase) UpdateOne(collection string, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()
	_, err := d.db.Collection(collection).UpdateOne(ctx, filter, update)
	return err
}

//method for upsert single value in a collection
func (d *mongoDatabase) UpsertOne(collection string, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeOutSecs))
	defer cancel()

	opts := options.Update().SetUpsert(true)
	_, err := d.db.Collection(collection).UpdateOne(ctx, filter, update, opts)
	return err
}
