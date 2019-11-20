package worker

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"log"
)

// CreateAccountPayload for call CreateAccount
type CreateAccountPayload struct {
	Username string
	Password string
	Email    string
}

// CreateAccount return a basic response
func CreateAccount(payload interface{}) interface{} {
	context := payload.(CreateAccountPayload)

	accounts, err := models.FindAllAccountByEmailOrUsername(context.Email, context.Username)
	if err != nil {
		log.Println("Error FindAccountByEmailOrUsername")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error server check account", -2)
	}

	if len(accounts) > 0 {
		errors := []string{}
		for _, account := range accounts {
			if helpers.StringCheckFormat(account.Username) == helpers.StringCheckFormat(context.Username) {
				errors = append(errors, "username")
			}
			if helpers.StringCheckFormat(account.Email) == helpers.StringCheckFormat(context.Email) {
				errors = append(errors, "email")
			}
		}
		return response.BasicResponse(errors, "Error account already exist", -3)
	}

	account, err := models.NewAccount(context.Username, context.Email, context.Password)
	if err != nil {
		log.Println("Error NewAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error server create account", -4)
	}

	return response.BasicResponse(account, "ok", 1)
}
