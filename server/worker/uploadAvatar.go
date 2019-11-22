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

// UploadAvatarPayload for call Upload
type UploadAvatarPayload struct {
	IDAccount int
	Avatar    []byte
	Ext       string
}

// UploadAvatarReturn it's the return of Upload
type UploadAvatarReturn struct {
	Avatars map[string]*models.Avatar
	Errors  []string
}

// CraftAvatarItem variable image
type CraftAvatarItem struct {
	size    uint
	quality int
}

var craftAvatar = map[string]CraftAvatarItem{
	"small":    CraftAvatarItem{size: 200, quality: 75},
	"original": CraftAvatarItem{size: 0, quality: 75},
}

// UploadAvatar return a basic response
func UploadAvatar(payload interface{}) interface{} {
	context := payload.(UploadAvatarPayload)

	// create a tmp file
	imgFileSource, err := ioutil.TempFile("images", "*"+context.Ext)
	if err != nil {
		log.Println("Error TempFile(imgFileSource) on UploadAvatar")
		log.Println(err)
		return response.BasicResponse(new(interface{}), "Error tmp file", -11)
	}
	// write into the tmp file the image bytes
	imgFileSource.Write(context.Avatar)
	imgFileSource.Close()

	var wg sync.WaitGroup
	uploadReturn := UploadAvatarReturn{Avatars: make(map[string]*models.Avatar), Errors: make([]string, 0)}

	// make a channel error for the goroutine
	errorsAvatar := make(chan error, len(craftAvatar))

	configuration, _ := config.GetConfig()

	for key, craftItem := range craftAvatar {
		wg.Add(1)
		go func(key string, size uint, quality int) {
			defer wg.Done()
			file, err := os.Open(imgFileSource.Name())
			if err != nil {
				log.Println("Error Open(file) on UploadAvatar")
				errorsAvatar <- err
				return
			}
			defer file.Close()

			imgFile, err := ioutil.TempFile("images", key+"-*"+context.Ext)
			if err != nil {
				log.Println("Error TempFile(imgFile) on UploadAvatar")
				errorsAvatar <- err
				return
			}
			defer imgFile.Close()

			imgFileBlured, err := ioutil.TempFile("images", key+"Blured-*"+context.Ext)
			if err != nil {
				log.Println("Error TempFile(imgFileBlured) on UploadAvatar")
				errorsAvatar <- err
				return
			}
			defer imgFileBlured.Close()

			pictureFile, _, err := image.Decode(file)
			if err != nil {
				log.Println("Error Decode(pictureFile) on UploadAvatar")
				errorsAvatar <- err
				return
			}

			picture, err := models.NewAvatar(context.IDAccount, configuration.BasePathImage+"/"+strings.ReplaceAll(imgFile.Name(), "images/", ""), configuration.BasePathImage+"/"+strings.ReplaceAll(imgFileBlured.Name(), "images/", ""), "server", key)
			if err != nil {
				log.Println("Error NewAvatar on UploadAvatar")
				errorsAvatar <- err
				return
			}

			uploadReturn.Avatars[key] = picture

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
		for err := range errorsAvatar {
			log.Println(err)
			uploadReturn.Errors = append(uploadReturn.Errors, err.Error())
		}
	}()

	wg.Wait()
	close(errorsAvatar)

	os.Remove(imgFileSource.Name())

	return response.BasicResponse(uploadReturn, "ok", 1)
}
