package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gastrader/hotelBE_go/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			// fmt.Println("Token not present in the header")
			return ErrUnauthortized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		//check token expiration

		expires := claims["expires"].(string)

		expiresTime, err := time.Parse(time.RFC3339Nano, expires)
		if err != nil {
			return fmt.Errorf("invalid bruh %v", err)
		}
		if !expiresTime.After(time.Now()) {
			return NewError(http.StatusUnauthorized, "token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnauthortized()
		}
		//set current authenticated user to the context..
		c.Context().SetUserValue("user", user)
		// fmt.Println("---> JWT auth")
		return c.Next()
	}
}
func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, ErrUnauthortized()
		}
		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token: ", err)
		return nil, ErrUnauthortized()
	}
	if !token.Valid {
		return nil, ErrUnauthortized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnauthortized()
	}
	return claims, nil

}
