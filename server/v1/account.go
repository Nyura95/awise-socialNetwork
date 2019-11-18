package v1

import (
	"encoding/json"
	"net/http"
)

// CreateAccount for create a new account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body login
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)
	json.NewEncoder(w).Encode("ok")
}
