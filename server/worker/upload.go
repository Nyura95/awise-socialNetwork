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
	Pictures map[string]*models.Picture
	Errors   []string
}

// CraftItem variable image
type CraftItem struct {
	size    uint
	radius  uint32
	quality int
}

var craft = map[string]CraftItem{
	"small-blured": CraftItem{radius: 50, size: 200, quality: 30},
	"small":        CraftItem{radius: 0, size: 200, quality: 75},
}

// Upload return a basic response
func Upload(payload interface{}) interface{} {
	context := payload.(UploadPayload)

	// create a tmp file
	imgFileSource, err := ioutil.TempFile("images", "original-*.jpg")
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
	uploadReturn := UploadReturn{Pictures: make(map[string]*models.Picture), Errors: make([]string, 0)}

	// make a channel error for the goroutine
	errorsPicture := make(chan error, len(craft))

	configuration, _ := config.GetConfig()

	for key, craftItem := range craft {
		wg.Add(1)
		go func(key string, size uint, radius uint32, quality int) {
			defer wg.Done()
			file, err := os.Open(imgFileSource.Name())
			if err != nil {
				errorsPicture <- err
				return
			}
			defer file.Close()
			imgFile, err := ioutil.TempFile("images", key+"-*.jpg")
			if err != nil {
				errorsPicture <- err
				return
			}
			defer imgFile.Close()

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

			uploadReturn.Pictures[key] = picture

			if radius > 0 {
				pictureFile = stackblur.Process(pictureFile, radius)
			}

			jpeg.Encode(imgFile, resize.Resize(size, 0, pictureFile, resize.Lanczos3), &jpeg.Options{Quality: quality})

		}(key, craftItem.size, craftItem.radius, craftItem.quality)
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

	return response.BasicResponse(uploadReturn, "ok", 1)
}
