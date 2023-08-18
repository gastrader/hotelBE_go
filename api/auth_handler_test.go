package api

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gastrader/hotelBE_go/db"
	"github.com/gastrader/hotelBE_go/types"
	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "Gavinito@tester.com",
		FirstName: "GAVY",
		LastName:  "PLUMO",
		Password:  "supersecret!!!!",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:     "Gavinito@tester.com",
		Password:  "supersecret!!!!",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP status of 200, but got %d", resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}
	if authResponse.Token == "" {
		t.Fatalf("expect JWT token to be present...")
	}
	//set encrypted password to empty string because ENCPW is not returned in any JSON response
	insertedUser.EncryptedPassword = ""
	// fmt.Println("Insert User is %s and auth res is %s", insertedUser.ID, authResponse.)
	if !reflect.DeepEqual(insertedUser, authResponse.User){
		t.Fatalf("expected user to be inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:     "Gavinito@tester.com",
		Password:  "supersecretWRONG",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected HTTP status of 400, but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected gen response type of be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid creds" {
		t.Fatalf("expected gen response type of be <invalid creds> but got %s", genResp.Msg)
	}
}