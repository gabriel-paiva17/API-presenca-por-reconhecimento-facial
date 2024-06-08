package main

import (
	"myproject/cv"
	"os"
	"fmt"
)


func main() {
	
	caminhoDoArquivo := "statics/2people_base64.txt"

    // Lendo o arquivo
    conteudo, err := os.ReadFile(caminhoDoArquivo)
    if err != nil {
		fmt.Printf("Erro ao ler o arquivo: %v", err)
		return
	}

    // Convertendo o conteúdo para string e atribuindo à variável
    base64image := string(conteudo)

	c, err := cv.CountFaces(base64image)

	fmt.Printf("Faces: %v\nErro: %v", c, err)

	// CreateServer()
}

