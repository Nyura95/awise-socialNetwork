package v1

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/server/response"
	"awise-socialNetwork/server/worker"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

// UploadPicture for create a new image
func UploadPicture(w http.ResponseWriter, r *http.Request) {
	IDUser := context.Get(r, "IDUser").(int)
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("picture")
	if err != nil {
		log.Printf("Error Retrieving image")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "Error retrieving image", -1))
		return
	}
	defer file.Close()

	pool := helpers.CreateWorkerPool(worker.UploadPicture)
	defer pool.Close()

	basicResponse := pool.Process(worker.UploadPicturePayload{IDAccount: IDUser, Picture: file}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}

// UploadAvatar for create a new image
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	IDUser := context.Get(r, "IDUser").(int)
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("Error Retrieving image")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "Error retrieving image", -1))
		return
	}
	defer file.Close()

	pool := helpers.CreateWorkerPool(worker.UploadAvatar)
	defer pool.Close()

	basicResponse := pool.Process(worker.UploadAvatarPayload{IDAccount: IDUser, Avatar: file}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}
