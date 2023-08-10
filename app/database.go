package app

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/dan13ram/wpokt-validator/models"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	lock "github.com/square/mongo-lock"
)

type Database interface {
	Connect() error
	SetupLockers() error
	SetupIndexes() error
	Disconnect() error
	InsertOne(collection string, data interface{}) error
	FindOne(collection string, filter interface{}, result interface{}) error
	FindMany(collection string, filter interface{}, result interface{}) error
	UpdateOne(collection string, filter interface{}, update interface{}) error
	UpsertOne(collection string, filter interface{}, update interface{}) error

	XLock(collection string, resourceId string) (string, error)
	SLock(collection string, resourceId string) (string, error)
	Unlock(collection string, lockId string) error
}

// mongoDatabase is a wrapper around the mongo database
type mongoDatabase struct {
	db *mongo.Database

	lockers map[string]*lock.Client
}

var (
	DB Database
)

// Connect connects to the database
func (d *mongoDatabase) Connect() error {
	log.Debug("[DB] Connecting to database")
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(time.Duration(Config.MongoDB.TimeoutSecs)*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Config.MongoDB.URI).SetWriteConcern(wcMajority))
	if err != nil {
		return err
	}
	d.db = client.Database(Config.MongoDB.Database)

	log.Info("[DB] Connected to mongo database: ", Config.MongoDB.Database)
	return nil
}

// SetupLocker sets up the locker
func (d *mongoDatabase) SetupLockers() error {
	log.Debug("[DB] Setting up locker")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	d.lockers = make(map[string]*lock.Client)
	var locker *lock.Client

	locker = lock.NewClient(d.db.Collection(models.CollectionMints))
	locker.CreateIndexes(ctx)
	d.lockers[models.CollectionMints] = locker

	locker = lock.NewClient(d.db.Collection(models.CollectionInvalidMints))
	locker.CreateIndexes(ctx)
	d.lockers[models.CollectionInvalidMints] = locker

	locker = lock.NewClient(d.db.Collection(models.CollectionBurns))
	locker.CreateIndexes(ctx)
	d.lockers[models.CollectionBurns] = locker

	locker = lock.NewClient(d.db.Collection(models.CollectionHealthChecks))
	locker.CreateIndexes(ctx)
	d.lockers[models.CollectionHealthChecks] = locker

	log.Info("[DB] Locker setup")
	return nil
}

func randomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// XLock locks a resource for exclusive access
func (d *mongoDatabase) XLock(collection string, resourceId string) (string, error) {
	locker := d.lockers[collection]
	if locker == nil {
		return "", errors.New("Locker not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()

	lockId := randomString(32)
	err := locker.XLock(ctx, resourceId, lockId, lock.LockDetails{})
	return lockId, err
}

// SLock locks a resource for shared access
func (d *mongoDatabase) SLock(collection string, resourceId string) (string, error) {
	locker := d.lockers[collection]
	if locker == nil {
		return "", errors.New("Locker not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()

	lockId := randomString(32)
	err := locker.SLock(ctx, resourceId, lockId, lock.LockDetails{}, -1)
	return lockId, err
}

// Unlock unlocks a resource
func (d *mongoDatabase) Unlock(collection string, lockId string) error {
	locker := d.lockers[collection]
	if locker == nil {
		return errors.New("Locker not found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()

	_, err := locker.Unlock(ctx, lockId)
	return err
}

// Setup Indexes
func (d *mongoDatabase) SetupIndexes() error {
	log.Debug("[DB] Setting up indexes")

	// setup unique index for mints
	log.Debug("[DB] Setting up indexes for mints")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	_, err := d.db.Collection(models.CollectionMints).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// setup unique index for invalid mints
	log.Debug("[DB] Setting up indexes for invalid mints")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	_, err = d.db.Collection(models.CollectionInvalidMints).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// setup unique index for burns
	log.Debug("[DB] Setting up indexes for burns")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	_, err = d.db.Collection(models.CollectionBurns).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_hash", Value: 1}, {Key: "log_index", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// setup unique index for healthchecks
	log.Debug("[DB] Setting up indexes for healthchecks")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	_, err = d.db.Collection(models.CollectionHealthChecks).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "validator_id", Value: 1}, {Key: "hostname", Value: 1}},
		Options: options.Index().SetUnique(true),
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	err := d.db.Client().Disconnect(ctx)
	log.Info("[DB] Disconnected from database")
	return err
}

// InitDB creates a new database wrapper
func InitDB() {
	DB = &mongoDatabase{}

	err := DB.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = DB.SetupIndexes()
	if err != nil {
		log.Fatal(err)
	}
	err = DB.SetupLockers()
	log.Info("[DB] Database initialized")
}

// method for insert single value in a collection
func (d *mongoDatabase) InsertOne(collection string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	_, err := d.db.Collection(collection).InsertOne(ctx, data)
	return err
}

// method for find single value in a collection
func (d *mongoDatabase) FindOne(collection string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	err := d.db.Collection(collection).FindOne(ctx, filter).Decode(result)
	return err
}

// method for find multiple values in a collection
func (d *mongoDatabase) FindMany(collection string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()
	_, err := d.db.Collection(collection).UpdateOne(ctx, filter, update)
	return err
}

//method for upsert single value in a collection
func (d *mongoDatabase) UpsertOne(collection string, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.MongoDB.TimeoutSecs))
	defer cancel()

	opts := options.Update().SetUpsert(true)
	_, err := d.db.Collection(collection).UpdateOne(ctx, filter, update, opts)
	return err
}
