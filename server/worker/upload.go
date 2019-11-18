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

	"github.com/esimov/stackblur-go"
	"github.com/nfnt/resize"
)

// UploadPayload for call Upload
type UploadPayload struct {
	Image multipart.File
}

// UploadReturn it's the return of Upload
type UploadReturn struct {
	Pictures []*models.Picture
	Errors   []string
}

var imgQuality = map[string][]int{
	"small":  []int{50},
	"medium": []int{75},
	"big":    []int{75},
}

var imgResize = map[string]uint{
	"small":  200,
	"medium": 500,
	"big":    1000,
}

// Upload return a basic response
func Upload(payload interface{}) interface{} {
	context := payload.(UploadPayload)

	// create a tmp file
	imgFileSource, err := ioutil.TempFile("images", "upload-*.jpg")
	if err != nil {
		log.Println("Error tmp file")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error tmp file", -2)
	}

	// read and get image bytes
	fileBytes, err := ioutil.ReadAll(context.Image)
	if err != nil {
		log.Println("Error read picture")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error read picture", -2)
	}

	// write into the tmp file the image bytes
	imgFileSource.Write(fileBytes)

	var wg sync.WaitGroup
	uploadReturn := UploadReturn{}

	errorsPicture := make(chan error, len(imgResize))

	configuration, _ := config.GetConfig()

	for key, size := range imgResize {
		wg.Add(1)
		go func(key string, size uint) {
			defer wg.Done()
			file, err := os.Open(imgFileSource.Name())
			if err != nil {
				errorsPicture <- err
				return
			}
			defer file.Close()
			imgFile, err := ioutil.TempFile("images", "upload-"+key+"-*.jpg")
			if err != nil {
				errorsPicture <- err
				return
			}
			defer imgFile.Close()
			draw.

			pictureFile, err := jpeg.Decode(file)
			if err != nil {
				errorsPicture <- err
				return
			}

			picture, err := models.NewPicture(configuration.BasePathImage+"/"+strings.ReplaceAll(imgFile.Name(), "images/", ""), "server", key)
			if err != nil {
				errorsPicture <- err
				return
			}
			uploadReturn.Pictures = append(uploadReturn.Pictures, picture)
			jpeg.Encode(imgFile, resize.Resize(size, 0, stackblur.Process(pictureFile, 60), resize.Lanczos3), &jpeg.Options{Quality: 30})

		}(key, size)
	}

	go func() {
		for err := range errorsPicture {
			log.Println("Error routine resize picture")
			log.Println(err)
			uploadReturn.Errors = append(uploadReturn.Errors, err.Error())
		}
	}()

	wg.Wait()
	close(errorsPicture)

	imgFileSource.Close()
	os.Remove(imgFileSource.Name())

	return response.BasicResponse(uploadReturn, "ok", 1)
}
