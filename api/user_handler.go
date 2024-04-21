package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jpmoraess/hotel-reservation/db"
	"github.com/jpmoraess/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		userID = c.Params("id")
		input  types.UpdateUserInput
	)
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	if err := h.userStore.UpdateUser(c.Context(), userID, &input); err != nil {
		return err
	}
	return c.JSON(input)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.SendStatus(204)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var input types.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	if err := input.Validate(c.Context()); err != nil {
		c.SendStatus(400)
		return c.JSON(err)
	}
	user, err := types.NewUserFromInput(input)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.SendStatus(404)
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
