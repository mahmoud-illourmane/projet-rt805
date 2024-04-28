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
	deviceResults := req.DeviceResults
	message := req.Message

	// Affiche les données reçues
	fmt.Println("Message reçu du client:", message)
	fmt.Println("Données reçues:", deviceResults)

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
