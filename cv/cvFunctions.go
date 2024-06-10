package cv

import (
	"encoding/base64"
	"fmt"
	"math"
    "strings"

	"github.com/Kagami/go-face"
)

const (
	modelDir = "faceModels"
)

/////////////////////////////////////////////////
// parece que só funciona com .jpg, multilinha //
/////////////////////////////////////////////////


// Decodifica uma imagem base64 e retorna os dados como []byte
func decodeBase64Image(base64Image string) ([]byte, error) {
	// Decodificar a imagem base64
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar imagem base64: %v", err)
	}

	return imageData, nil
}

// Conta o número de rostos em uma imagem base64
func CountFaces(base64Image string) (int, error) {
	recognizer, err := face.NewRecognizer(modelDir)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar reconhecedor: %v", err)
	}
	defer recognizer.Close()

	imgData, err := decodeBase64Image(base64Image)
	if err != nil {
		return 0, err
	}

	faces, err := recognizer.Recognize(imgData)
	if err != nil {
		return 0, fmt.Errorf("erro ao reconhecer rostos na imagem: %v", err)
	}

	return len(faces), nil
}

// Calcula a distância euclidiana entre dois descritores faciais
func euclideanDistance(descriptor1, descriptor2 face.Descriptor) float32 {
	var sum float32
	for i := range descriptor1 {
		diff := descriptor1[i] - descriptor2[i]
		sum += diff * diff
	}
	return float32(math.Sqrt(float64(sum)))
}

// Compara dois rostos em imagens base64 para ver se são da mesma pessoa
func CompareFaces(base64Image1, base64Image2 string) (bool, error) {
	recognizer, err := face.NewRecognizer(modelDir)
	if err != nil {
		return false, fmt.Errorf("erro ao criar reconhecedor: %v", err)
	}
	defer recognizer.Close()

	imgData1, err := decodeBase64Image(base64Image1)
	if err != nil {
		return false, err
	}

	imgData2, err := decodeBase64Image(base64Image2)
	if err != nil {
		return false, err
	}

	faces1, err := recognizer.Recognize(imgData1)
	if err != nil {
		return false, fmt.Errorf("erro ao reconhecer rostos na primeira imagem: %v", err)
	}

	faces2, err := recognizer.Recognize(imgData2)
	if err != nil {
		return false, fmt.Errorf("erro ao reconhecer rostos na segunda imagem: %v", err)
	}

	if len(faces1) == 0 || len(faces2) == 0 {
		return false, fmt.Errorf("nenhum rosto encontrado em uma das imagens")
	}

	face1 := faces1[0].Descriptor
	face2 := faces2[0].Descriptor

	distance := euclideanDistance(face1, face2)

	const threshold = 0.6
	isSamePerson := distance < threshold

	return isSamePerson, nil
}

func IsBase64JPG(s string) bool {
	// Verifica se a string começa com o prefixo base64 de uma imagem JPEG
	if !strings.HasPrefix(s, "/9j/") {
		return false
	}

	// Decodifica a string base64
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}

	// Verifica o cabeçalho JPEG (primeiros 3 bytes: 0xFF, 0xD8, 0xFF)
	if len(data) < 3 || data[0] != 0xFF || data[1] != 0xD8 || data[2] != 0xFF {
		return false
	}

	return true
}