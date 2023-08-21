package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gastrader/hotelBE_go/db/fixtures"
	"github.com/gastrader/hotelBE_go/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "bar hotel", "montreal", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 98.99, hotel.ID)
	booking := fixtures.AddBooking(db.Store, room.ID, user.ID, time.Now(), time.Now().AddDate(0,0,2), 2)
	_ = booking

	app := fiber.New()
	bookingHandler := NewBookingHandler(db.Store)
	app.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookings)
}