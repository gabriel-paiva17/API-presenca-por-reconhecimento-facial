package main

import("fmt"
		)

func main() {
	
	imagepath1 := "imgs/face_detected_20240523-160423.png"
	imagepath2, err := DetectAndSaveFace()

	if err != nil {

		return

	}

	facesIguais, err := AreFacesEqual(imagepath1, imagepath2)

	if err != nil {

		fmt.Printf("erro inesperado")
		return

	}


	if facesIguais {

		fmt.Printf("\nigual\nnao vai ser salvo")
		DeleteFile(imagepath2)
		return

	}

	fmt.Printf("\ndiferente\nsalvo!\n")

}

