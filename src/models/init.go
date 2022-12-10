package models

import (
	"fmt"
	"gosfV2/src/models/env"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

var (
	UserIconDir = filepath.Join(env.Config.FilesDirectory, "icons")
	DefaultIcon = filepath.Join(UserIconDir, "default.jpeg")
)

func init() {
	if err := os.MkdirAll(env.Config.FilesDirectory, 0744); err != nil {
		fmt.Println("Error creating files directory: ", err.Error())
		os.Exit(1)
	}

	if err := os.MkdirAll(UserIconDir, 0744); err != nil {
		fmt.Println("Error creating icon files directory: ", err.Error())
		os.Exit(1)
	}

	file, err := os.Create(DefaultIcon)
	if err != nil {
		fmt.Println("Error creating default icon: ", err.Error())
		os.Exit(1)
	}

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	img.Set(0, 0, image.White)

	if err := jpeg.Encode(file, img, nil); err != nil {
		fmt.Println("Error creating default icon: ", err.Error())
		os.Exit(1)
	}

	if err := file.Close(); err != nil {
		fmt.Println("Error creating default icon: ", err.Error())
		os.Exit(1)
	}
}
