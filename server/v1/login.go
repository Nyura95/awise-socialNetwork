package v1

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/server/response"
	"awise-socialNetwork/server/worker"
	"encoding/json"
	"log"
	"net/http"
)

type login struct {
	username string
	password string
}

// Login authenticate an user
func Login(w http.ResponseWriter, r *http.Request) {

	var body login
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	if body.username == "" || body.password == "" {
		log.Printf("Body login invalid")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "The body for login is not valid", -1))
		return
	}

	pool := helpers.CreateWorkerPool(worker.Login)
	defer pool.Close()

	basicResponse := pool.Process(worker.LoginPayload{Password: body.password, Username: body.username}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}
