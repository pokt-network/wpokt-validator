package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRemoveDuplicates(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()

	elements := []primitive.ObjectID{id1, id2, id1, id3, id2, id3}
	expected := []primitive.ObjectID{id1, id2, id3}

	result := RemoveDuplicates(elements)
	assert.Equal(t, expected, result)
}

func TestRemoveDuplicates_NoDuplicates(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()

	elements := []primitive.ObjectID{id1, id2, id3}
	expected := []primitive.ObjectID{id1, id2, id3}

	result := RemoveDuplicates(elements)
	assert.Equal(t, expected, result)
}

func TestRemoveDuplicates_EmptySlice(t *testing.T) {
	elements := []primitive.ObjectID{}
	expected := []primitive.ObjectID{}

	result := RemoveDuplicates(elements)
	assert.Equal(t, expected, result)
}
