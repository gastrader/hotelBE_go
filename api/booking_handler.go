package api

import (
	"fmt"
	"github.com/gastrader/hotelBE_go/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil{
		return ErrResourceNotFound("booking")
	}
	fmt.Printf("C CONTEXT IS: %+v", c.Context().UserValue("user"))

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthortized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthortized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"cancelled": true}); err != nil{
		return err
	}
	return c.JSON(genericResp{Type: "msg", Msg: "updated"})
}


func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

//TODO: this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthortized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthortized()
	}
	return c.JSON(booking)
}