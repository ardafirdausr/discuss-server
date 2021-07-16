package mongo

import (
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toObjectID(ID interface{}) (primitive.ObjectID, error) {
	strID, ok := ID.(string)
	if !ok {
		err := entity.ErrNotFound{Message: "Invalid ID"}
		return primitive.NilObjectID, err
	}

	objID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		err = entity.ErrNotFound{
			Message: "Failed get data using the corresponding ID",
			Err:     err,
		}
		return primitive.NilObjectID, err
	}

	return objID, nil
}
