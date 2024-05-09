package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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

/*==== Structures de MongoDB ====*/

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

/*==== END/Structures de MongoDB ====*/

/*
*	Cette méthode établi la connexion avec la BD.
*	Elle prend en paramètre l'url de connexion
 */
func NewMongoDBClient(uri string) (*MongoDBClient, error) {
	// Crée une instance de client MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	// Etabli une connexion
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Vérification de la connexion
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connecté à MongoDB")
	return &MongoDBClient{client: client}, nil
}

/*
*	Cette méthode permet de fermer la connexion à MongoDB.
 */
func (c *MongoDBClient) Close() error {
	// Fermer la connexion MongoDB
	return c.client.Disconnect(context.Background())
}

/*
* Fonction qui affiche les données de manière
* indifférente en fonction du nom du dispositif
 */
func (c *MongoDBClient) GetDataByDeviceName(databaseName, collectionName, deviceName string) error {
	collection := c.client.Database(databaseName).Collection(collectionName)

	filter := bson.D{{Key: "device_name", Value: deviceName}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	fmt.Printf("Données pour le dispositif '%s' :\n", deviceName)
	for cursor.Next(context.Background()) {
		var data Data
		if err := cursor.Decode(&data); err != nil {
			return err
		}
		fmt.Printf("Journée : %d, SuccessCounts : %v, FailureCounts : %v\n", data.Day, data.SuccessCounts, data.FailureCounts)
	}
	if cursor.Err() != nil {
		return cursor.Err()
	}
	return nil
}

/*
*	Cette méthode permet d'établir la connexion
* 	et met à jour les données existantes ou insère de nouvelles données.
 */
func (c *MongoDBClient) addDataToMongoDB(deviceResults *pb.DeviceResults, journee int32, deviceName string) {
	// Création des données à insérer ou à mettre à jour
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

	// Filtre
	filter := bson.M{
		"deviceName": deviceName,
	}

	// Mettre à jour les champs SuccessCounts et FailureCounts
	update := bson.M{
		"$set": bson.M{
			"successCounts": data.SuccessCounts,
			"failureCounts": data.FailureCounts,
		},
	}

	// Options pour upsert
	options := options.Update().SetUpsert(true)

	// Mise à jour ou insertion (upsert)
	_, err := c.client.Database("projet-805").Collection("devices_data").UpdateOne(context.Background(), filter, update, options)

	// Vérification des erreurs
	if err != nil {
		log.Fatalf("Erreur lors de la mise à jour ou de l'insertion des données : %v", err)
	}
}

/*
* Cette fonction affiche toutes les données
* qui se trouve en BD
**/
func (c *MongoDBClient) GetAllData(databaseName, collectionName string) error {
	// Accès à la collection
	collection := c.client.Database(databaseName).Collection(collectionName)

	// Récuperation de toutes les données de la collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Erreur lors de la recherche des données : %v", err)
		return err
	}
	defer cursor.Close(context.Background())

	fmt.Println("Données extraites :")

	// Parcour les documents renvoyés par le curseur
	for cursor.Next(context.Background()) {
		var data struct {
			DeviceName    string `bson:"deviceName"`
			SuccessCounts []struct {
				Key   string `bson:"key"`
				Value int    `bson:"value"`
			} `bson:"successCounts"`
			FailureCounts []struct {
				Key   string `bson:"key"`
				Value int    `bson:"value"`
			} `bson:"failureCounts"`
		}

		// Décode les données dans la variable data
		err = cursor.Decode(&data)
		if err != nil {
			log.Printf("Erreur lors du décodage des données : %v", err)
			return err
		}

		// Affiche les données extraites
		fmt.Printf("Nom du dispositif : %s\n", data.DeviceName)
		fmt.Println("SuccessCounts :")
		for _, count := range data.SuccessCounts {
			fmt.Printf("    %s : %d\n", count.Key, count.Value)
		}
		fmt.Println("FailureCounts :")
		for _, count := range data.FailureCounts {
			fmt.Printf("    %s : %d\n", count.Key, count.Value)
		}
		fmt.Println()
	}

	// Vérifie si le curseur a rencontré une erreur lors de l'itération
	if err = cursor.Err(); err != nil {
		log.Printf("Erreur lors de l'itération des données : %v", err)
		return err
	}

	return nil
}
