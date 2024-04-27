package main

import (
	"log"
)

func main() {
	filePath := "../donnees/journee_1.json"

	results, err := ParseFile(filePath)
	if err != nil {
		log.Fatalf("erreur lors de l'analyse du fichier : %v", err)
	}

	DisplayResults(results)
}
