package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jpmoraess/hotel-reservation/db"
	"github.com/jpmoraess/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (testdb *testdb) teatDown(t *testing.T) {
	if err := testdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TESTDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	testdb := setup()
	defer testdb.teatDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	input := types.CreateUserInput{
		FirstName: "John",
		LastName:  "Wick",
		Email:     "john_wick@mail.com",
		Password:  "j0hnw1ck",
	}

	b, _ := json.Marshal(input)

	request := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	request.Header.Add("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(response.Body).Decode(&user)

	if response.StatusCode != 200 {
		t.Fail()
	}
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the response")
	}
	if user.FirstName != input.FirstName {
		t.Errorf("expected firstname %s but got %s", input.FirstName, user.FirstName)
	}
	if user.LastName != input.LastName {
		t.Errorf("expected lastname %s but got %s", input.LastName, user.LastName)
	}
	if user.Email != input.Email {
		t.Errorf("expected email %s but got %s", input.Email, user.Email)
	}
}
