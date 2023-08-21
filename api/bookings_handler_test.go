package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gastrader/hotelBE_go/db/fixtures"
	"github.com/gastrader/hotelBE_go/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser = fixtures.AddUser(db.Store, "jimbo", "nonad", false)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "montreal", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 98.99, hotel.ID)
		booking   = fixtures.AddBooking(db.Store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2)
		app       = fiber.New()
		route = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	_ = booking
	
	
	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code, got %d", resp.StatusCode)
	}
	var bookingReponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingReponse); err != nil {
		t.Fatal(err)
	}
	have := bookingReponse
	if have.ID != booking.ID {
		t.Fatal("expected matching IDs")
	}
	if have.UserID != booking.UserID {
		t.Fatal("expected matching user IDs")
	}
	
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("non 200 code, got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "montreal", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 98.99, hotel.ID)
		booking   = fixtures.AddBooking(db.Store, room.ID, user.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2)
		app       = fiber.New()
	)
	_ = booking
	admin := app.Group("/", JWTAuthentication(db.User), AdminAuth)
	bookingHandler := NewBookingHandler(db.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatal("expected matching IDs")
	}
	if have.UserID != booking.UserID {
		t.Fatal("expected matching user IDs")
	}

	// test non-admin cannot access bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf(" expected non 200 response but got %d", resp.StatusCode)
	}
}
