package middleware

import (
	"time"

	"github.com/AdluAghnia/not_todolist/database"
	"github.com/AdluAghnia/not_todolist/models"
	"github.com/AdluAghnia/not_todolist/repository"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
)

var jwtSecretKey = []byte("SECRET_KEY")

func GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
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

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid JWT",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid JWT Claims",
			})
		}

		userID := claims["id"].(float64)
		c.Locals("userID", userID)

		return c.Next()
	}
}

func GetUserFromContext(c *fiber.Ctx) (*models.User, error) {
	db, err := database.Db()
	if err != nil {
		return nil, err
	}
	userID := c.Locals("userID").(float64)
	user, err := repository.GetUserByID(db, uint(userID))

	if err != nil {
		return nil, err
	}

	return user, nil
}
