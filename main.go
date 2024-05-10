package main

import (
	"context"
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/jpmoraess/hotel-reservation/api"
	"github.com/jpmoraess/hotel-reservation/api/middleware"
	"github.com/jpmoraess/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}

	// handlers initialization
	var (
		userStore  = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)

		store = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}

		authHandler  = api.NewAuthHandler(userStore)
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		roomHandler  = api.NewRoomHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiV1        = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// booking handler
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	app.Listen(*listenAddr)
}
