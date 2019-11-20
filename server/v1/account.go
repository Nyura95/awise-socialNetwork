package v1

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/server/response"
	"awise-socialNetwork/server/worker"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type createAccount struct {
	Username string
	Password string
	Email    string
}

// CreateAccount for create a new account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body createAccount
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	if body.Username == "" || body.Password == "" || body.Email == "" {
		log.Printf("Body createAccount invalid")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "The body for createAccount is not valid (need username and password)", -1))
		return
	}

	pool := helpers.CreateWorkerPool(worker.CreateAccount)
	defer pool.Close()

	basicResponse := pool.Process(worker.CreateAccountPayload{Username: body.Username, Password: body.Password, Email: body.Email}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}

// GetAccount get account info
func GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("GetAccount parsing query params error")
		log.Println(err)
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "id query is not valid", -1))
		return
	}

	pool := helpers.CreateWorkerPool(worker.GetAccount)
	defer pool.Close()

	basicResponse := pool.Process(worker.GetAccountPayload{ID: id}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}
