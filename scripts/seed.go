package main

import (
	"context"
	"log"

	"github.com/jpmoraess/hotel-reservation/db"
	"github.com/jpmoraess/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	roomStore := db.NewMongoRoomStore(client, db.DBNAME)
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "France",
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 88.99,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 122.99,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.99,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	insertedRoomsId := []primitive.ObjectID{}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		insertedRoomsId = append(insertedRoomsId, insertedRoom.ID)
	}

	if err := hotelStore.Update(ctx, insertedHotel.ID.Hex(), insertedRoomsId); err != nil {
		log.Fatal(err)
	}
}
