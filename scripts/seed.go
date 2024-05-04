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

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func main() {
	seedUser("John", "Wick", "john_wick@mail.com", "johnwick123")

	seedHotel("Bellucia", "France", 5)
	seedHotel("The Cozy Hotel", "The Nederlands", 4)
	seedHotel("Copacabana Palace", "Rio de Janeiro", 5)
}

func seedUser(firstName string, lastName string, email string, password string) error {
	user, err := types.NewUserFromInput(types.CreateUserInput{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func seedHotel(name string, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 88.99,
		},
		{
			Size:  "normal",
			Price: 122.99,
		},
		{
			Size:  "kingsize",
			Price: 199.99,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	userStore = db.NewMongoUserStore(client)
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
