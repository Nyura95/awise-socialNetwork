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

// UploadPicturePayload for call Upload
type UploadPicturePayload struct {
	IDAccount int
	Picture   multipart.File
}

// UploadPictureReturn it's the return of Upload
type UploadPictureReturn struct {
	Pictures map[string]*models.Picture
	Errors   []string
}

// CraftPictureItem variable image
type CraftPictureItem struct {
	size    uint
	radius  uint32
	quality int
}

var craftPicture = map[string]CraftPictureItem{
	"small-blured":    CraftPictureItem{radius: 50, size: 200, quality: 30},
	"small":           CraftPictureItem{radius: 0, size: 200, quality: 75},
	"original":        CraftPictureItem{radius: 0, size: 0, quality: 75},
	"original-blured": CraftPictureItem{radius: 50, size: 0, quality: 30},
}

// UploadPicture return a basic response
func UploadPicture(payload interface{}) interface{} {
	context := payload.(UploadPicturePayload)

	// create a tmp file
	imgFileSource, err := ioutil.TempFile("images", "*.jpg")
	if err != nil {
		log.Println("Error tmp file")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error tmp file", -2)
	}

	// read and get image bytes
	fileBytes, err := ioutil.ReadAll(context.Picture)
	if err != nil {
		log.Println("Error read picture")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error read picture", -2)
	}

	// write into the tmp file the image bytes
	imgFileSource.Write(fileBytes)

	var wg sync.WaitGroup
	uploadReturn := UploadPictureReturn{Pictures: make(map[string]*models.Picture), Errors: make([]string, 0)}

	// make a channel error for the goroutine
	errorsPicture := make(chan error, len(craftPicture))

	configuration, _ := config.GetConfig()

	for key, craftItem := range craftPicture {
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

			imgFileBlured, err := ioutil.TempFile("images", key+"Blured-*.jpg")
			if err != nil {
				errorsPicture <- err
				return
			}
			defer imgFileBlured.Close()

			pictureFile, err := jpeg.Decode(file)
			if err != nil {
				errorsPicture <- err
				return
			}

			picture, err := models.NewPicture(context.IDAccount, configuration.BasePathImage+"/"+strings.ReplaceAll(imgFile.Name(), "images/", ""), configuration.BasePathImage+"/"+strings.ReplaceAll(imgFileBlured.Name(), "images/", ""), "server", key)
			if err != nil {
				errorsPicture <- err
				return
			}

			uploadReturn.Pictures[key] = picture

			jpeg.Encode(imgFile, resize.Resize(size, 0, pictureFile, resize.Lanczos3), &jpeg.Options{Quality: quality})
			jpeg.Encode(imgFileBlured, resize.Resize(size, 0, stackblur.Process(pictureFile, 50), resize.Lanczos3), &jpeg.Options{Quality: 30})

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
	os.Remove(imgFileSource.Name())

	return response.BasicResponse(uploadReturn, "ok", 1)
}
