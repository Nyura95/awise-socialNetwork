package v1

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/server/response"
	"awise-socialNetwork/server/worker"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/context"
)

// UploadPicture for create a new image
func UploadPicture(w http.ResponseWriter, r *http.Request) {
	IDUser := context.Get(r, "IDUser").(int)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("picture")
	if err != nil {
		log.Printf("Body uploadPicture invalid")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "The body for uploadPicture is not valid", -1))
		return
	}
	defer file.Close()

	bFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error read file")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "Error reading image", -2))
		return
	}

	pool := helpers.CreateWorkerPool(worker.UploadPicture)
	defer pool.Close()

	basicResponse := pool.Process(worker.UploadPicturePayload{IDAccount: IDUser, Picture: bFile, Ext: filepath.Ext(handler.Filename)}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}

// UploadAvatar for create a new image
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	IDUser := context.Get(r, "IDUser").(int)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("Body uploadAvatar invalid")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "The body for uploadAvatar is not valid", -1))
		return
	}
	defer file.Close()

	bFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error read file")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "Error reading image", -2))
		return
	}

	pool := helpers.CreateWorkerPool(worker.UploadAvatar)
	defer pool.Close()

	basicResponse := pool.Process(worker.UploadAvatarPayload{IDAccount: IDUser, Avatar: bFile, Ext: filepath.Ext(handler.Filename)}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}
