package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gastrader/hotelBE_go/api"
	"github.com/gastrader/hotelBE_go/db"
	"github.com/gastrader/hotelBE_go/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
	}

	user := fixtures.AddUser(store, "Amanda", "Joyce", false)
	fmt.Printf("%s -> %s\n", user.FirstName, api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "Admin", "Strator", true)
	fmt.Printf("%s -> %s\n", admin.FirstName, api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 199.99, hotel.ID)
	booking := fixtures.AddBooking(store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 3), 2)
	fmt.Printf("Test booking -> %+v \n", booking)
	fmt.Println("-----> DB Seeded <-----")

	for i := 0; i<100; i++ {
		name := fmt.Sprintf("Random %d", i)
		fixtures.AddHotel(store, name, "Canada", rand.Intn(5)+1, nil)
	}
}
