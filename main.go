package main

import("fmt")

func main() {
	
	imagepath1 := "imgs/face_detected_20240523-130527.png"
	imagepath2, err := DetectAndSaveFace()

	if err != nil {

		return

	}

	facesIguais, err := AreFacesEqual(imagepath1, imagepath2)

	if err != nil {

		fmt.Printf("erro inesperado")

	}


	if facesIguais {

		fmt.Printf("igual")
		return

	}

	fmt.Printf("diferente")

}

