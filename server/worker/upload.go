package worker

import (
	"awise-socialNetwork/config"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

// UploadPayload for call Upload
type UploadPayload struct {
	Image multipart.File
}

var imgResize = map[string]uint{
	"small":  200,
	"medium": 500,
	"big":    1000,
}

// Upload return a basic response
func Upload(payload interface{}) interface{} {
	context := payload.(UploadPayload)

	imgFileSource, err := ioutil.TempFile("images", "upload-*.jpg")
	if err != nil {
		log.Println("Error tmp file")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error tmp file", -2)
	}

	fileBytes, err := ioutil.ReadAll(context.Image)
	if err != nil {
		log.Println("Error read picture")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error read picture", -2)
	}

	imgFileSource.Write(fileBytes)

	var wg sync.WaitGroup
	var pictures []*models.Picture

	errors := make(chan error, len(imgResize))
	defer close(errors)

	configuration, _ := config.GetConfig()

	for key, size := range imgResize {
		wg.Add(1)
		go func(key string, size uint) {
			file, err := os.Open(imgFileSource.Name())
			if err != nil {
				errors <- err
			}
			defer file.Close()
			imgFile, err := ioutil.TempFile("images", "upload-"+key+"-*.jpg")
			if err != nil {
				errors <- err
			}
			defer func() {
				imgFile.Close()
				wg.Done()
			}()
			pictureFile, err := jpeg.Decode(file)
			if err != nil {
				errors <- err
			}
			picture, err := models.NewPicture(configuration.BasePathImage+"/"+strings.ReplaceAll(imgFile.Name(), "images/", ""), "server", key)
			if err != nil {
				errors <- err
			}
			pictures = append(pictures, picture)
			jpeg.Encode(imgFile, resize.Resize(size, 0, pictureFile, resize.Lanczos3), nil)
		}(key, size)
	}

	wg.Wait()

	if len(errors) != 0 {
		log.Println("Error create pictures")
		for err := range errors {
			log.Println(err)
		}
	}

	imgFileSource.Close()
	os.Remove(imgFileSource.Name())

	return response.BasicResponse(pictures, "ok", 1)
}
