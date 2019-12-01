package worker

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// GoogleLoginPayload for call GoogleLogin
type GoogleLoginPayload struct {
	Token string
}

// GoogleResponse is a struct of google response
type GoogleResponse struct {
	ErrorDescription string `json:"error_description"`
	Email            string `json:"email"`
	EmailVerified    bool   `json:"email_verified"`
	Name             string `json:"name"`
	Picture          string `json:"Picture"`
	GivenName        string `json:"given_name"`
	FamilyName       string `json:"family_name"`
}

// GoogleLogin return a basic response
func GoogleLogin(payload interface{}) interface{} {
	context := payload.(GoogleLoginPayload)

	res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + context.Token)
	if err != nil {
		log.Println("Error http on GoogleLogin")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find account avatars", -11)
	}
	defer res.Body.Close()

	var body GoogleResponse
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&body)

	if body.ErrorDescription != "" && body.Email == "" {
		log.Println("Error body on GoogleLogin")
		return response.BasicResponse(body, "Error find account", -11)
	}

	account, err := models.FindAccountByEmail(body.Email)
	if err != nil {
		log.Println("Error FindAccountByEmail on GoogleLogin")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find account account", -11)
	}

	if account.ID == 0 {
		account, err = models.NewAccount(strings.Split(body.Email, "@")[0], body.Email, "")
	}

	accessToken, err := models.NewAccessToken(account.ID, helpers.GetUtc().AddDate(0, 1, 0))
	if err != nil {
		log.Println("Error NewAccessToken on Login")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create token", -11)
	}

	return response.BasicResponse(accessToken, "ok", 1)
}
