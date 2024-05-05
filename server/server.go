package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "mahmoud.projet.rt0805/proto/SendData"
)

// Server est la structure qui implémente l'interface SendData définie dans votre fichier `.proto`.
type Server struct {
	pb.UnimplementedSendDataServer
}

// RpcSendData est la méthode qui reçoit les données envoyées par le client.
func (s *Server) RpcSendData(ctx context.Context, req *pb.SendDataRequest) (*pb.SendDataReply, error) {
	// Reception des données de la requête
	// deviceResults := req.DeviceResults
	// journee := req.Journee
	// deviceName := req.DeviceResults.DeviceName

	// Extraire success_count CREATE de la requête
	successCountCREATE := req.DeviceResults.SuccessCount

	// Rechercher la clé "CREATE" dans success_count
	createValue, exists := successCountCREATE["CREATE"]
	if exists {
		fmt.Printf("success_count CREATE : {key:\"CREATE\" value:%d}\n", createValue)
	} else {
		fmt.Println("La clé 'CREATE' n'existe pas dans success_count\n")
	}

	// fmt.Println("Device Name : \n", deviceName)

	// // Affiche les données reçues
	// fmt.Println("Données reçu :\nJournée :", journee)
	// fmt.Println("Données :\n", deviceResults)

	addDataToMongoDB()
	// TODO

	// Réponse au client
	return &pb.SendDataReply{
		Message: "Données reçues avec succès",
	}, nil
}

func main() {
	// Adresse et le port sur lesquels le serveur écoutera
	address := "localhost:50051"

	// Listener TCP sur l'adresse et le port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("échec de l'écoute sur %s : %v", address, err)
	}

	// Création du serveur gRPC
	grpcServer := grpc.NewServer()

	// Enregistre le service `SendData` sur le serveur
	pb.RegisterSendDataServer(grpcServer, &Server{})

	// Démarre le serveur et écoute les connexions entrantes
	log.Printf("Serveur gRPC en écoute sur %s...", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("échec du démarrage du serveur : %v", err)
	}
}
