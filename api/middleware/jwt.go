package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("---> JWT auth")
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized my guy")
	}
	fmt.Println("THe token is:", token)

	claims, err := validateToken(token); 
	if err != nil {
		return err
	}
	
	//check token expiration

	expires := claims["expires"].(string)

	expiresTime, err := time.Parse(time.RFC3339Nano, expires)
	if err != nil {
		return fmt.Errorf("invalid bruh")
	}
	if !expiresTime.After(time.Now()) {
		return fmt.Errorf("invalid time")
	}
	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token: ", err)
		return nil, fmt.Errorf("unauthorized dog")
	}
	if !token.Valid{
		return nil, fmt.Errorf("token not valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil

}
