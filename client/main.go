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

	// Extractions des données depuis le fichier.
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

	// Créeation d'un tableau de pointeurs
	var resultsPointers []*pb.DeviceResults
	for i := range results {
		resultsPointers = append(resultsPointers, &results[i])
	}

	// Envoi de la requête gRPC
	response, erreur := client.RpcSendData(ctx, &pb.SendDataRequest{
		DeviceResults: resultsPointers,
		Journee:       int32(journee),
	})
	if erreur != nil {
		log.Fatalf("Échec de l'envoi des données : %v", erreur)
	}
	log.Printf("Réponse de la part du serveur : %s", response.GetMessage())
}
