package main

import (
	"bytes"

	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
	"fmt"
	"os"
	"image"
	"image/jpeg"
)

type CaptchaController struct {
	beego.Controller
}



func main() {
	// Create a new Captcha
	ValidationString := captcha.New()

	// Store the string for validation later
	StoreString := ValidationString

	fmt.Println(StoreString)
	// Create the captcha image
	var ImageBuffer bytes.Buffer
	captcha.WriteImage(&ImageBuffer, ValidationString, 300, 90)

	img,_,_ := image.Decode(bytes.NewReader(ImageBuffer.Bytes()))
	outputFile, err := os.Create("src/github.com/geobe/gostip/resources/images/test.jpg")
	if err != nil {
		panic(error(err))
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	jpeg.Encode(outputFile, img,nil)
	fmt.Println(captcha.Reload(StoreString))

	outputFile.Close()
}