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
		log.Println("Error FindAccountByEmailOrUsername on CreateAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find account", -11)
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
		return response.BasicResponse(errors, "Account already exist", -13)
	}

	account, err := models.NewAccount(context.Username, context.Email, context.Password)
	if err != nil {
		log.Println("Error NewAccount on CreateAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create account", -11)
	}

	Avatar, err := models.NewAvatar(account.ID, "https://image.awise.co/s9igio6ytzc885ngio10i.jpg", "https://image.awise.co/s9igio6ytzc885ngio10i.jpg", "default", "default")
	if err != nil {
		log.Println("Error NewAvatar on CreateAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create account", -11)
	}

	return response.BasicResponse(models.AccountWithAvatar{Account: account, Avatar: Avatar}, "ok", 1)
}
