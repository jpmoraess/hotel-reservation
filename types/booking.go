package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	RoomID      primitive.ObjectID `bson:"roomID" json:"roomID"`
	NumOfPerson int                `bson:"numOfPersons" json:"numOfPersons"`
	FromDate    time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate    time.Time          `bson:"tillDate" json:"tillDate"`
}
