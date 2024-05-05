package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	pb "mahmoud.projet.rt0805/proto/SendData"
)

/*
*
*	Classe qui me pert d'établir une connexion à la base de donnée MongoDB.
*	Elle se charge également de définir les fonctions d'intéractions
*	avec la base de donnée.
*
 */

/* Structures de MongoDB */

type MongoDBClient struct {
	client *mongo.Client
}

type OperationCount struct {
	Key   string `bson:"key"`   // CREATE UPDATE DELETE
	Value int    `bson:"value"` // Nb operations
}

type Data struct {
	Day           int              `bson:"day"`
	DeviceName    string           `bson:"device_name"`
	SuccessCounts []OperationCount `bson:"success_counts"`
	FailureCounts []OperationCount `bson:"failure_counts"`
}

/*
*	Cette méthode établi la connexion avec la BD.
*	Elle prend en paramètre l'url de connexion
 */
func NewMongoDBClient(uri string) (*MongoDBClient, error) {
	// Créer une instance de client MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	// Etabli une connexion
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Vérifier la connexion
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connecté à MongoDB")
	return &MongoDBClient{client: client}, nil
}

/*
*	Cette méthode d'insérer les données.
 */
func (c *MongoDBClient) InsertData(databaseName, collectionName string, data Data) error {
	// Obtenir la collection
	collection := c.client.Database(databaseName).Collection(collectionName)

	// Insérer les données
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	fmt.Println("Données insérées avec succès")
	return nil
}

/*
*	Cette méthode permet de fermer la connexion.
 */
func (c *MongoDBClient) Close() error {
	// Fermer la connexion MongoDB
	return c.client.Disconnect(context.Background())
}

/*
*	Cette méthode permet d'établir la connexion
* 	et l'ajout de données en BD.
 */
func addDataToMongoDB(deviceResults *pb.DeviceResults, journee int32, deviceName string) {
	uri := "mongodb://root:root@10.22.9.96:27017/"

	// Créer un client MongoDB
	client, erreur := NewMongoDBClient(uri)
	if erreur != nil {
		log.Fatalf("Erreur lors de la création du client MongoDB : %v", erreur)
	}
	defer client.Close()

	// Créer des données à insérer
	data := Data{
		Day:        int(journee),
		DeviceName: deviceName,
		SuccessCounts: []OperationCount{
			{Key: "CREATE", Value: int(deviceResults.SuccessCount["CREATE"])},
			{Key: "DELETE", Value: int(deviceResults.SuccessCount["DELETE"])},
			{Key: "UPDATE", Value: int(deviceResults.SuccessCount["UPDATE"])},
		},
		FailureCounts: []OperationCount{
			{Key: "CREATE", Value: int(deviceResults.FailureCount["CREATE"])},
			{Key: "DELETE", Value: int(deviceResults.FailureCount["DELETE"])},
			{Key: "UPDATE", Value: int(deviceResults.FailureCount["UPDATE"])},
		},
	}

	// Insérer les données dans la base de données
	erreur = client.InsertData("projet-805", "devices_data", data)
	if erreur != nil {
		log.Fatalf("Erreur lors de l'insertion des données : %v", erreur)
	}
}
