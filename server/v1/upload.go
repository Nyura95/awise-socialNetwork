package v1

import (
	"awise-socialNetwork/helpers"
	"awise-socialNetwork/server/response"
	"awise-socialNetwork/server/worker"
	"encoding/json"
	"log"
	"net/http"
)

// UploadImage for create a new image
func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error Retrieving image")
		json.NewEncoder(w).Encode(response.BasicResponse(new(interface{}), "Error retrieving image", -1))
		return
	}
	defer file.Close()

	pool := helpers.CreateWorkerPool(worker.Upload)
	defer pool.Close()

	basicResponse := pool.Process(worker.UploadPayload{Image: file}).(response.Response)
	if basicResponse.Success == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(basicResponse)
}
