package v1

import (
	"encoding/json"
	"net/http"
)

// CreateAccount for create a new account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("ok")
}
