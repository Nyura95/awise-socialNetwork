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
		log.Println("Error FindAccount on GetAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find account", -11)
	}

	if account.ID == 0 {
		log.Println("Account not found on GetAccount")
		return response.BasicResponse(new(interface{}), "Account not found", -12)
	}

	accountAvatar, err := models.FindAccountAvatarByIDAccountActive(account.ID)
	if err != nil {
		log.Println("Error FindAllAccountAvatarsByIDAccount on GetAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find account avatars", -11)
	}

	if accountAvatar.ID == 0 {
		return response.BasicResponse(models.AccountWithAvatar{Account: account}, "ok", 1)
	}

	avatar, err := models.FindAvatar(accountAvatar.ID)
	if err != nil {
		log.Println("Error FindAvatar on GetAccount")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find avatar", -11)
	}

	return response.BasicResponse(models.AccountWithAvatar{Account: account, Avatar: avatar}, "ok", 1)
}
