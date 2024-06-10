package main

/*
import (	
	"fmt"
	"myproject/cv"
	"os"
)
*/

func main() {
	
	CreateServer()
	
	// testing CV functions
	/*
	caminhoDoArquivo1 := "statics/lJames.txt"
	caminhoDoArquivo2 := "statics/messi2.txt"


	// Lendo o arquivo
	conteudo1, err := os.ReadFile(caminhoDoArquivo1)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo: %v", err)
		return
	}

	conteudo2, err := os.ReadFile(caminhoDoArquivo2)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo: %v", err)
		return
	}

	// Convertendo o conteúdo para string e atribuindo à variável
	base64image1 := string(conteudo1)
	base64image2 := string(conteudo2)

	c, err := cv.CountFaces(base64image1)

	fmt.Printf("Faces: %v\nErro: %v\n\n", c, err)

	// CreateServer()

	r, err := cv.CompareFaces(base64image1, base64image2)

	fmt.Printf("Resultado: %v,\nErro: %v", r, err)
	*/
}
