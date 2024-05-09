package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "mahmoud.projet.rt0805/proto/SendData"
)

// Server est la structure qui implémente l'interface SendData définie dans votre fichier `.proto`.
type Server struct {
	pb.UnimplementedSendDataServer
	mongoClient *MongoDBClient
}

// Constructeur pour Server
func NewServer(client *MongoDBClient) *Server {
	return &Server{
		mongoClient: client,
	}
}

/*
* RpcSendData est la méthode qui reçoit les données envoyées par le client.
 */
func (s *Server) RpcSendData(ctx context.Context, req *pb.SendDataRequest) (*pb.SendDataReply, error) {
	if req.DeviceResults == nil {
		log.Fatalf("Erreur : Aucune donnée n'est disponible.")
		// Réponse au client
		return &pb.SendDataReply{
			Message: "Erreur : Aucune donnée n'est disponible.",
		}, nil
	}

	verif := VerifyData(req)
	if !verif {
		log.Fatalf("Erreur : Les données reçues ne sont pas correctes.")
		// Réponse au client
		return &pb.SendDataReply{
			Message: "Erreur : Les données reçues ne sont pas correctes.",
		}, nil
	}

	for _, deviceResult := range req.DeviceResults {
		// Vérifie que deviceResult n'est pas nil pour éviter des erreurs de référence nil
		if deviceResult != nil {
			deviceName := deviceResult.DeviceName

			s.mongoClient.addDataToMongoDB(deviceResult, req.Journee, deviceName)
		}
	}

	/*== Affichage des résultats ==*/
	err := s.mongoClient.GetAllData("projet-805", "devices_data")
	if err != nil {
		log.Fatalf("Echec GetDataByDeviceName : %v", err)
	}

	// Réponse au client
	return &pb.SendDataReply{
		Message: "Données reçues avec succès",
	}, nil
}

func main() {
	/*== CONNEXION A LA BD ==*/
	uri := "mongodb://root:root@10.22.9.96:27017/"
	// Création d'un client MongoDB
	client, erreur := NewMongoDBClient(uri)
	if erreur != nil {
		log.Fatalf("Erreur lors de la création du client MongoDB : %v", erreur)
	}
	// Cette méthode assure que la fonction close sera appelé juste avant que la fonction se termine.
	defer client.Close()

	/*== SERVER gRPC ==*/
	address := "localhost:50051" // Adresse et le port sur lesquels le serveur écoutera

	// Listener TCP sur l'adresse et le port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("échec de l'écoute sur %s : %v", address, err)
	}

	// Création du serveur gRPC
	grpcServer := grpc.NewServer()

	// Création d'une instance de Server avec le client MongoDB
	server := NewServer(client)

	// Enregistre le service `SendData` sur le serveur
	pb.RegisterSendDataServer(grpcServer, server)

	// Démarre le serveur et écoute les connexions entrantes
	log.Printf("Serveur gRPC en écoute sur %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("échec du démarrage du serveur : %v", err)
	}
}
