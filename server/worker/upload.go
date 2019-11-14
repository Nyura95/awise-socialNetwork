package worker

import (
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

// UploadPayload for call Upload
type UploadPayload struct {
	Image multipart.File
}

// Upload return a basic response
func Upload(payload interface{}) interface{} {
	context := payload.(UploadPayload)

	imgFile, err := ioutil.TempFile("images", "upload-*.png")
	if err != nil {
		log.Println("Error tmp file")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error tmp file", -2)
	}
	defer imgFile.Close()

	fileBytes, err := ioutil.ReadAll(context.Image)
	if err != nil {
		log.Println("Error read picture")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error read picture", -2)
	}

	imgFile.Write(fileBytes)

	picture, err := models.NewPicture(strings.ReplaceAll(imgFile.Name(), "images/", ""), "server")
	if err != nil {
		log.Println("Error create picture")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error create picture", -3)
	}

	go func(path string) {
		file, _ := os.Open(path)
		img, _ := jpeg.Decode(file)
		file.Close()
		m := resize.Resize(200, 0, img, resize.Lanczos3)
		imgFile, _ := ioutil.TempFile("images", "upload-*.png")
		defer imgFile.Close()
		jpeg.Encode(imgFile, m, nil)

	}(imgFile.Name())

	return response.BasicResponse(picture, "ok", 1)
}
