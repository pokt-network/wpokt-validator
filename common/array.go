package common

import "go.mongodb.org/mongo-driver/bson/primitive"

func RemoveDuplicates(elements []primitive.ObjectID) []primitive.ObjectID {
	encountered := map[string]bool{}
	result := []primitive.ObjectID{}

	for _, element := range elements {
		if !encountered[element.String()] {
			encountered[element.String()] = true
			result = append(result, element)
		}
	}
	return result
}
