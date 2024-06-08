package cv

import (
    "encoding/base64"
    "fmt"
 
    "gocv.io/x/gocv"
)

func CountFaces(base64Image string) (int, error) {
    // Decodificar a string base64 para um array de bytes
    imageData, err := base64.StdEncoding.DecodeString(base64Image)
    if err != nil {
        return 0, fmt.Errorf("erro ao decodificar imagem base64: %v", err)
    }

    // Ler a imagem a partir dos bytes decodificados
    img, err := gocv.IMDecode(imageData, gocv.IMReadColor)
    if err != nil {
        return 0, fmt.Errorf("erro ao ler imagem: %v", err)
    }
    defer img.Close()

    // Carregar o classificador de detecção de faces
    classifier := gocv.NewCascadeClassifier()
    defer classifier.Close()

    if !classifier.Load("statics/haarcascade_frontalface_default.xml") {
        return 0, fmt.Errorf("erro ao carregar classificador de faces")
    }

    // Detectar rostos na imagem
    rects := classifier.DetectMultiScale(img)

    return len(rects), nil
}