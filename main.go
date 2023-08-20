package main

import (
	"context"
	"flag"
	"log"

	"github.com/gastrader/hotelBE_go/api"
	"github.com/gastrader/hotelBE_go/api/middleware"
	"github.com/gastrader/hotelBE_go/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenaddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	// handler init
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store      = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
			Booking: bookingStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler = api.NewAuthHandler(userStore)
		roomHandler = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("/api/")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	)

	//auth handler
	auth.Post("/auth", authHandler.HandleAuthenticate)

	//Versioned API routes
	//user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//room handlers
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	
	//TODO: Cancel booking!

	//booking handlers
	apiv1.Get("/booking", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	app.Listen(*listenaddr)

}
