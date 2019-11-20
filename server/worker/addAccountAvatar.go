package worker

import (
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"log"
)

// AddAccountAvatarPayload for call AddAccountAvatar
type AddAccountAvatarPayload struct {
	IDAccount int
	IDAvatar  int
}

// AddAccountAvatar return a basic response
func AddAccountAvatar(payload interface{}) interface{} {
	context := payload.(AddAccountAvatarPayload)

	avatar, err := models.FindAvatar(context.IDAvatar)
	if err != nil {
		log.Println("Error FindAvatar")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find avatar id", -2)
	}

	if avatar.ID == 0 {
		log.Println("Avatar not found")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "This avatar does not exist", -2)
	}

	accountAvatar, err := models.NewAccountAvatar(context.IDAccount, context.IDAvatar)
	if err != nil {
		log.Println("Error NewAccountAvatar")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create account avatar", -2)
	}

	return response.BasicResponse(accountAvatar, "ok", 1)
}
