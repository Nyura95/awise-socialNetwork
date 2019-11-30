package v1

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/server/response"
	"awise-socialNetwork/server/worker"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshLogin struct {
	RefreshToken string `json:"refresh_token"`
}

// Login authenticate an user
func Login(w http.ResponseWriter, r *http.Request) {

	var body login
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	if body.Email == "" || body.Password == "" {
		log.Printf("Body login invalid")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "The body for login is not valid", -1))
		return
	}

	pool := helpers.CreateWorkerPool(worker.Login)
	defer pool.Close()

	basicResponse := pool.Process(worker.LoginPayload{Password: body.Password, Email: strings.ToLower(body.Email)}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}

// RefreshLogin authenticate an user with the refresh token
func RefreshLogin(w http.ResponseWriter, r *http.Request) {

	var body refreshLogin
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	if body.RefreshToken == "" {
		log.Printf("Body refresh login invalid")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "The body for refresh login is not valid", -1))
		return
	}

	pool := helpers.CreateWorkerPool(worker.RefreshLogin)
	defer pool.Close()

	basicResponse := pool.Process(worker.RefreshLoginPayload{RefreshToken: body.RefreshToken}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)

}
