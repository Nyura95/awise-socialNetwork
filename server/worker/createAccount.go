package worker

import (
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
)

// CreateAccountPayload for call CreateAccount
type CreateAccountPayload struct {
	IDAvatars []int
	Username  string
	FirstName string
	LastName  string
	Password  string
	IDScope   int
	Score     int
	Level     int
	Credits   int
	Phone     string
	City      string
	Country   string
}

// CreateAccount return a basic response
func CreateAccount(payload interface{}) interface{} {
	context := payload.(CreateAccountPayload)

	models.NewAccount(context.IDAvatars, context.FirstName, context.LastName, context.Username, context.Score, context.Level, context.Credits, context.Phone, context.City, context.Country, context.Password, context.IDScope)

	return response.BasicResponse(context, "ok", 1)
}
