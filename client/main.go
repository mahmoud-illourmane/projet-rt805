package main

import (
	"context"
	"log"
	"strconv"
	"time"

	pb "mahmoud.projet.rt0805/proto/SendData"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Nom du fichier à traiter
	journee := 1
	fileName := "journee_" + strconv.Itoa(journee) + ".json"
	filePath := "../donnees/" + fileName

	results, erreur := ParseFile(filePath)
	if erreur != nil {
		log.Fatalf("erreur lors de l'analyse du fichier : %v", erreur)
	}

	// Établir une connexion avec le serveur.
	connexion, erreur := grpc.Dial(address, grpc.WithInsecure())
	if erreur != nil {
		log.Fatalf("Échec de la connexion : %v", erreur)
	}
	defer connexion.Close()
	client := pb.NewSendDataClient(connexion)

	// Contact avec le serveur et recevoir une réponse.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Convertir les résultats en un pointeur
	resultsPointer := &results[0] // Envoi du device numéro 1
	response, erreur := client.RpcSendData(ctx, &pb.SendDataRequest{
		DeviceResults: resultsPointer,
		Journee:       1,
	})
	if erreur != nil {
		log.Fatalf("Échec de l'envoi des données : %v", erreur)
	}
	log.Printf("Réponse de la part du serveur : %s", response.GetMessage())

	// DisplayResults(results)
}
