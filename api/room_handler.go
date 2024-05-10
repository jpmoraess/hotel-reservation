package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jpmoraess/hotel-reservation/db"
	"github.com/jpmoraess/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

// 2023-10-10T20:20:05Z
type BookRoomInput struct {
	FromDate     time.Time `json:"fromDate"`
	TillDate     time.Time `json:"tillDate"`
	NumOfPersons int       `json:"numOfPersons"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	var input BookRoomInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		c.SendStatus(http.StatusInternalServerError)
	}
	booking := types.Booking{
		UserID:      user.ID,
		RoomID:      roomID,
		FromDate:    input.FromDate,
		TillDate:    input.TillDate,
		NumOfPerson: input.NumOfPersons,
	}
	fmt.Printf("%+v", booking)
	return nil
}
