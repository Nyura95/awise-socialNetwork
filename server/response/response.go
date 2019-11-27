package response

import (
	"awise-socialNetwork/helpers"
	"os"
	"time"
)

// Response generic
type Response struct {
	StatusCode int         `json:"statusCode"`
	Reason     int         `json:"reason"`
	Comment    string      `json:"comment"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Hash       string      `json:"hash"`
	Version    string      `json:"version"`
	ServerTime time.Time   `json:"serverTime"`
}

// BasicResponse from API
func BasicResponse(data interface{}, comment string, reason int) Response {
	success := false
	statusCode := 400
	if reason == 1 {
		success = true
		statusCode = 200
	}
	basicResponse := Response{StatusCode: statusCode, Reason: reason, Comment: comment, Success: success, Data: data, Hash: helpers.StringToMD5(helpers.Stringify(data)), ServerTime: helpers.GetUtc(), Version: os.Getenv("VERSION")}
	return basicResponse
}
