package util

import "go.mongodb.org/mongo-driver/bson/primitive"

func StringIDToObjectID(ids []string) ([]primitive.ObjectID, err) {
	var err error
	objectIDs := make([]primitive.ObjectID, len(ids))

	for i, id := range ids {
		objectIDs[i], err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
	}
	return objectIDs, nil
}
