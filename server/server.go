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
	if req.DeviceResults == nil {
		log.Fatalf("Erreur : aucune donnée n'est accessible")
	}

	for _, deviceResult := range req.DeviceResults {
		// Vérifie que deviceResult n'est pas nil pour éviter des erreurs de référence nil
		if deviceResult != nil {
			deviceName := deviceResult.DeviceName
			successCount := deviceResult.SuccessCount
			failureCount := deviceResult.FailureCount

			addDataToMongoDB(deviceResult, req.Journee, deviceName)

			// Affichages
			fmt.Printf("\nNom de l'appareil : %s\n", deviceName)
			fmt.Printf("Nombre de succès : %v\n", successCount)
			fmt.Printf("Nombre d'échecs : %v\n", failureCount)
			fmt.Printf("\n\n")
		}
	}

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
	log.Printf("Serveur gRPC en écoute sur %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("échec du démarrage du serveur : %v", err)
	}
}
