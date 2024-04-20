package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jpmoraess/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "John",
		LastName:  "Wick",
	}

	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {

	return c.JSON("John")
}