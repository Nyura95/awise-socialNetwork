package worker

import (
	"awise-socialNetwork/config"
	"awise-socialNetwork/models"
	"awise-socialNetwork/server/response"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/esimov/stackblur-go"
	"github.com/nfnt/resize"
)

// UploadPicturePayload for call Upload
type UploadPicturePayload struct {
	IDAccount int
	Picture   []byte
	Ext       string
}

// UploadPictureReturn it's the return of Upload
type UploadPictureReturn struct {
	Pictures map[string]*models.Picture
	Errors   []string
}

// CraftPictureItem variable image
type CraftPictureItem struct {
	size    uint
	quality int
}

var craftPicture = map[string]CraftPictureItem{
	"small":    CraftPictureItem{size: 200, quality: 75},
	"original": CraftPictureItem{size: 0, quality: 75},
}

// UploadPicture return a basic response
func UploadPicture(payload interface{}) interface{} {
	context := payload.(UploadPicturePayload)

	// create a tmp file
	imgFileSource, err := ioutil.TempFile("images", "*"+context.Ext)
	if err != nil {
		log.Println("Error TempFile(imgFileSource) on UploadPicture")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error tmp file", -11)
	}
	// write into the tmp file the image bytes
	imgFileSource.Write(context.Picture)
	imgFileSource.Close()

	var wg sync.WaitGroup
	uploadReturn := UploadPictureReturn{Pictures: make(map[string]*models.Picture), Errors: make([]string, 0)}

	// make a channel error for the goroutine
	errorsPicture := make(chan error, len(craftPicture))

	configuration, _ := config.GetConfig()

	for key, craftItem := range craftPicture {
		wg.Add(1)
		go func(key string, size uint, quality int) {
			defer wg.Done()
			file, err := os.Open(imgFileSource.Name())
			if err != nil {
				log.Println("Error Open(file) on UploadPicture")
				errorsPicture <- err
				return
			}
			defer file.Close()

			imgFile, err := ioutil.TempFile("images", key+"-*"+context.Ext)
			if err != nil {
				log.Println("Error TempFile(imgFile) on UploadPicture")
				errorsPicture <- err
				return
			}
			defer imgFile.Close()

			imgFileBlured, err := ioutil.TempFile("images", key+"Blured-*"+context.Ext)
			if err != nil {
				log.Println("Error TempFile(imgFileBlured) on UploadPicture")
				errorsPicture <- err
				return
			}
			defer imgFileBlured.Close()

			pictureFile, _, err := image.Decode(file)
			if err != nil {
				log.Println("Error Decode(pictureFile) on UploadPicture")
				errorsPicture <- err
				return
			}

			picture, err := models.NewPicture(context.IDAccount, configuration.BasePathImage+"/"+strings.ReplaceAll(imgFile.Name(), "images/", ""), configuration.BasePathImage+"/"+strings.ReplaceAll(imgFileBlured.Name(), "images/", ""), "server", key)
			if err != nil {
				log.Println("Error NewPicture on UploadPicture")
				errorsPicture <- err
				return
			}

			uploadReturn.Pictures[key] = picture

			switch context.Ext {
			case ".png":
				// enc := &png.Encoder{
				// 	CompressionLevel: png.DefaultCompression,
				// }
				png.Encode(imgFile, resize.Resize(size, 0, pictureFile, resize.Lanczos3))
				png.Encode(imgFileBlured, resize.Resize(size, 0, stackblur.Process(pictureFile, 50), resize.Lanczos3))
			case ".jpg":
				jpeg.Encode(imgFile, resize.Resize(size, 0, pictureFile, resize.Lanczos3), &jpeg.Options{Quality: quality})
				jpeg.Encode(imgFileBlured, resize.Resize(size, 0, stackblur.Process(pictureFile, 50), resize.Lanczos3), &jpeg.Options{Quality: 30})
			}

		}(key, craftItem.size, craftItem.quality)
	}

	go func() {
		for err := range errorsPicture {
			log.Println(err)
			uploadReturn.Errors = append(uploadReturn.Errors, err.Error())
		}
	}()

	wg.Wait()
	close(errorsPicture)

	os.Remove(imgFileSource.Name())

	return response.BasicResponse(uploadReturn, "ok", 1)
}
