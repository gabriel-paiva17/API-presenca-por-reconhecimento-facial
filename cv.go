package main

import (
	"fmt"
	"image/color"
	"time"
	"strings"
	"os/exec"
	"os"

	"gocv.io/x/gocv"
)

func DetectAndSaveFace() (string, error) {
	// set to use a video capture device 0
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return "", err
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
	if !classifier.Load("./static/haarcascade_frontalface_default.xml") { 
		err := "error reading cascade file: data/haarcascade_frontalface_default.xml"
		fmt.Println(err)
		return "", fmt.Errorf(err)
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)

	filename := ""

	for {
		if ok := webcam.Read(&img); !ok {
			err := fmt.Sprintf("cannot read device %v\n", deviceID)
			fmt.Println(err)
			return "", fmt.Errorf(err)
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
			filename = SaveImage(img, key)
		}

		if filename != "" {

			return filename, nil

		}

	}
}

// caso aperte S ou s, a imagem é salva 
func SaveImage(img gocv.Mat, key int) string {

	filename := ""

	if key == 'S' || key == 's' {

		timestamp := time.Now().Format("20060102-150405")
		filename = fmt.Sprintf("imgs/face_detected_%s.png", timestamp)
		
		if ok := gocv.IMWrite(filename, img); ok {
			fmt.Printf("Image saved to %s\n", filename)
			time.Sleep(3 * time.Second)
			
		}
		
	}

	return filename

}

// CompareFaces executa um script Python para comparar faces em duas imagens e retorna um valor booleano.
func AreFacesEqual(imagePath1, imagePath2 string) (bool, error) {
	// Comando para executar o script Python com os caminhos das imagens como argumentos
	cmd := exec.Command("python3", "compare_faces.py", imagePath1, imagePath2)

	// Executa o comando e captura a saída
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error executing python script: %v", err)
	}

	// Converte a saída para uma string e remove espaços em branco extras
	result := strings.TrimSpace(string(output))

	// Verifica a saída e retorna true ou false conforme apropriado
	return result == "true", nil

}

func DeleteFile(filePath string) error {

	// Tenta deletar o arquivo
	err := os.Remove(filePath)

	if err != nil {
		fmt.Printf("Error deleting file: %v\n", err)
		return err
	}

	fmt.Println("File deleted successfully")
	return nil

}
