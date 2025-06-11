package cosmos

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"

	log "github.com/sirupsen/logrus"
)

type resultMaxSequence struct {
	MaxSequence uint64 `bson:"max_sequence"`
}

func findMaxSequenceFromInvalidMints() (*uint64, error) {
	filter := bson.M{
		"sequence":      bson.M{"$ne": nil},
		"vault_address": app.Config.Pocket.MultisigAddress,
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "max_sequence", Value: bson.D{{Key: "$max", Value: "$sequence"}}},
		}}},
	}

	var result resultMaxSequence
	err := app.DB.AggregateOne(models.CollectionInvalidMints, pipeline, &result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	maxSequence := uint64(result.MaxSequence)

	return &maxSequence, nil
}

func findMaxSequenceFromBurns() (*uint64, error) {
	filter := bson.M{
		"sequence":      bson.M{"$ne": nil},
		"vault_address": app.Config.Pocket.MultisigAddress,
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "max_sequence", Value: bson.D{{Key: "$max", Value: "$sequence"}}},
		}}},
	}

	var result resultMaxSequence
	err := app.DB.AggregateOne(models.CollectionBurns, pipeline, &result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	maxSequence := uint64(result.MaxSequence)

	return &maxSequence, nil
}

func findMaxSequence() (*uint64, error) {
	maxSequenceInvalidMints, err := findMaxSequenceFromInvalidMints()
	if err != nil {
		return nil, err
	}

	maxSequenceBurns, err := findMaxSequenceFromBurns()
	if err != nil {
		return nil, err
	}

	if maxSequenceInvalidMints == nil && maxSequenceBurns == nil {
		return nil, nil
	}

	if maxSequenceInvalidMints == nil {
		return maxSequenceBurns, nil
	}

	if maxSequenceBurns == nil {
		return maxSequenceInvalidMints, nil
	}

	if *maxSequenceInvalidMints > *maxSequenceBurns {
		return maxSequenceInvalidMints, nil
	}

	return maxSequenceBurns, nil
}

var FindMaxSequence = findMaxSequence

const sequenceResourseID = "comsos_sequence"

func lockReadSequences() (lockID string, err error) {
	lockID, err = app.DB.SLock(sequenceResourseID)
	if err != nil {
		log.WithError(err).Error("Error locking max sequence")
		return
	}
	log.WithField("resource_id", sequenceResourseID).Debug("Locked read sequences")
	return
}

var LockReadSequences = lockReadSequences

func lockWriteSequence() (lockID string, err error) {
	lockID, err = app.DB.XLock(sequenceResourseID)
	if err != nil {
		log.WithError(err).Error("Error locking max sequence")
		return
	}
	log.WithField("resource_id", sequenceResourseID).Debug("Locked write sequence")
	return
}

var LockWriteSequence = lockWriteSequence
