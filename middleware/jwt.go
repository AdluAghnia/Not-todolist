package middleware

import (
	"time"

	"github.com/AdluAghnia/not_todolist/models"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
)

var jwtSecretKey = []byte("SECRET_KEY")

func GenerateJWT(user *models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
        "id": user.ID,
        "email": user.Email,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })

    return token.SignedString(jwtSecretKey)
}

func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Cookies("jwt")

        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Missing or invalid JWT",
            })
        }

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            return jwtSecretKey, nil
        })

        if err != nil || !token.Valid{
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid JWT",
            })
        }

        return c.Next()
    }
}
