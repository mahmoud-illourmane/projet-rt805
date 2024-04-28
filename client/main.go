package main

import (
	"context"
	"log"
	"time"

	pb "mahmoud.projet.rt0805/proto/SendData"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	filePath := "../donnees/journee_1.json"

	results, err := ParseFile(filePath)
	if err != nil {
		log.Fatalf("erreur lors de l'analyse du fichier : %v", err)
	}

	// Établir une connexion avec le serveur.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Échec de la connexion : %v", err)
	}
	defer conn.Close()
	client := pb.NewSendDataClient(conn)

	// Contact avec le serveur et recevoir une réponse.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Convertir les résultats en un pointeur vers une tranche
	resultsPointer := &results[0]
	response, err := client.RpcSendData(ctx, &pb.SendDataRequest{
		DeviceResults: resultsPointer,
		Message:       "Envoi de résultats.",
	})

	if err != nil {
		log.Fatalf("Échec de l'appel : %v", err)
	}
	log.Printf("Réponse : %s", response.GetMessage())

	// DisplayResults(results)
}
