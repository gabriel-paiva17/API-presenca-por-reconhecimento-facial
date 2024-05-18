package main

import (
	"fmt"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

func DetectFaces() {
	// set to use a video capture device 0
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	// path para haar cascade da biblioteca
	if !classifier.Load("/home/gabriel/src-go/gocv/data/haarcascade_frontalface_default.xml") { 
		fmt.Println("Error reading cascade file: data/haarcascade_frontalface_default.xml")
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 1)
		}

		window.IMShow(img)

		key := window.WaitKey(1)

		// se só há uma face detectada, permite o salvamento da imagem
		if len(rects) == 1 {
			
			fmt.Println("Press 's' to save the image")
			SaveImage(img, key)
			
		}
	}
}

// caso aperte S ou s, a imagem é salva 
func SaveImage(img gocv.Mat, key int) {

	if key == 'S' || key == 's' {

		timestamp := time.Now().Format("20060102-150405")
		filename := fmt.Sprintf("imgs/face_detected_%s.png", timestamp)
		
		if ok := gocv.IMWrite(filename, img); ok {
			fmt.Printf("Image saved to %s\n", filename)
			time.Sleep(3 * time.Second)
			
		}
		
	}

}

