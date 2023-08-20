package middleware

import (
	"fmt"

	"github.com/gastrader/hotelBE_go/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not auth'd")
	}
	if !user.IsAdmin {
		return fmt.Errorf("not auth'd")
	}
	return c.Next()
}