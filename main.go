package main

import (
	"context"
	"log"
	"os"

	"github.com/gastrader/hotelBE_go/api"
	"github.com/gastrader/hotelBE_go/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CONFIG
// 1. Mongo endpoint
// 2. Listen Address of HTTP server
// 3. JWT secret
// 4. MongoDBName

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	// listenaddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	// flag.Parse()

	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	// handler init
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api/")

		apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", api.AdminAuth)
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
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	//admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	app.Listen(listenAddr)

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
