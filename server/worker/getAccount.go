package worker

import (
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"log"
)

// GetAccountPayload for call GetAccount
type GetAccountPayload struct {
	ID int
}

// GetAccount return a basic response
func GetAccount(payload interface{}) interface{} {
	context := payload.(GetAccountPayload)

	account, err := models.FindAccount(context.ID)
	if err != nil {
		log.Println("Error FindAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error server find account", -2)
	}

	if account.ID == 0 {
		log.Println("Account not found")
		return response.BasicResponse(new(interface{}), "Account not found", -3)
	}

	return response.BasicResponse(account, "ok", 1)
}
