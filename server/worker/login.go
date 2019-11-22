package worker

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"log"
)

// LoginPayload for call Login
type LoginPayload struct {
	Email    string
	Password string
}

// Login return a basic response
func Login(payload interface{}) interface{} {
	context := payload.(LoginPayload)

	account, err := models.FindAccountByPassword(context.Email, context.Password)
	if err != nil {
		log.Println("Error FindAccountByPassword on Login")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find account", -11)
	}

	if account.ID == 0 || account.Email != context.Email {
		log.Println("Password or username not valid on Login")
		return response.BasicResponse(new(interface{}), "Password or username is not valid", -12)
	}

	// err = models.DeleteAllAccessTokenByIDAccount(account.ID)
	// if err != nil {
	// 	log.Println("Error, disabled all token")
	// 	log.Println(err)
	// 	return response.BasicResponse(new(interface{}), "Error disabled token", -4)
	// }

	accessToken, err := models.NewAccessToken(account.ID, helpers.GetUtc().AddDate(0, 1, 0))
	if err != nil {
		log.Println("Error NewAccessToken on Login")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create token", -11)
	}

	return response.BasicResponse(accessToken, "ok", 1)
}
