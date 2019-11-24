package worker

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"log"
)

// RefreshLoginPayload for call RefreshLogin
type RefreshLoginPayload struct {
	RefreshToken string
}

// RefreshLogin return a basic response
func RefreshLogin(payload interface{}) interface{} {
	context := payload.(RefreshLoginPayload)

	accessToken, err := models.FindAccessTokenByRefreshToken(context.RefreshToken)
	if err != nil {
		log.Println("Error FindAccessTokenByRefreshToken on RefreshLogin")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error find access token", -11)
	}

	if accessToken.ID == 0 {
		log.Println("accessToken not found on RefreshLogin")
		return response.BasicResponse(new(interface{}), "Refresh token does not exist", -12)
	}

	newAccessToken, err := models.NewAccessToken(accessToken.IDAccount, helpers.GetUtc().AddDate(0, 1, 0))
	if err != nil {
		log.Println("Error NewAccessToken on RefreshLogin")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create token", -11)
	}

	return response.BasicResponse(newAccessToken, "ok", 1)
}
