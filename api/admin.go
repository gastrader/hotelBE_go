package api

import (
	"github.com/gastrader/hotelBE_go/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthortized()
	}
	if !user.IsAdmin {
		return ErrUnauthortized()
	}
	return c.Next()
}