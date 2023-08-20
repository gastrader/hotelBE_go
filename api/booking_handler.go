package api

import (
	"github.com/gastrader/hotelBE_go/db"
	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

//TODO: this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	return nil
}

//TODO: this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	return nil
}