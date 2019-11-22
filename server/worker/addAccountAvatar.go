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
		log.Println("Error FindAvatar on AddAccountAvatar worker")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find avatar id", -11)
	}

	if avatar.ID == 0 {
		log.Println("Avatar not found on AddAccountAvatar worker")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "This avatar does not exist", -12)
	}

	err = models.DisabledAllAvatarByIDAccount(context.IDAccount)
	if err != nil {
		log.Println("Error DisabledAllAvatarByIDAccount on AddAccountAvatar")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error delete avatar", -2)
	}

	accountAvatar, err := models.NewAccountAvatar(context.IDAccount, context.IDAvatar)
	if err != nil {
		log.Println("Error NewAccountAvatar on AddAccountAvatar")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create account avatar", -2)
	}

	return response.BasicResponse(accountAvatar, "ok", 1)
}
