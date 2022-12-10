package models

import (
	"fmt"
	"gosfV2/src/models/env"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

var (
	UserIconDir = filepath.Join(env.Config.FilesDirectory, "icons")
	DefaultIcon = filepath.Join(UserIconDir, "default.jpeg")
)

func init() {
	if err := os.MkdirAll(env.Config.FilesDirectory, 0744); err != nil {
		log.Fatal("Error creating files directory: ", err.Error())
	}

	if err := os.MkdirAll(UserIconDir, 0744); err != nil {
		log.Fatal("Error creating icon files directory: ", err.Error())
	}

	file, err := os.Create(DefaultIcon)
	if err != nil {
		fmt.Println("Error creating default icon: ", err.Error())
		os.Exit(1)
	}

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	img.Set(0, 0, image.White)

	if err := jpeg.Encode(file, img, nil); err != nil {
		log.Fatal("Error creating default icon: ", err.Error())
	}

	if err := file.Close(); err != nil {
		log.Fatal("Error creating default icon: ", err.Error())
	}
}
