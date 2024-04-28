package main

/*
 *
 *	Ce fichier contient toutes les fonctions utiles à la partie
 *	client du projet.
 *	Auteur: Mahmoud Illourmane
 *	Date : 27/04/2024
 *
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	pb "mahmoud.projet.rt0805/proto/SendData"
)

/*== Public Functions  ==*/

/*
*	Fonction qui analyse le fichier JSON et extrait les résultats
*	Exemple d'un appel depuis /client :
*		filePath := "../donnees/journee_1.json"
*		results, err := ParseFile(filePath)
 */
func ParseFile(filePath string) ([]pb.DeviceResults, error) {
	var results []pb.DeviceResults

	// Ouvrir le fichier JSON
	file, err := os.Open(filePath)
	if err != nil {
		return results, fmt.Errorf("erreur d'ouverture du fichier : %v", err)
	}
	defer file.Close()

	// Lecture du fichier
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return results, fmt.Errorf("erreur de lecture du fichier : %v", err)
	}

	// Décoder le contenu JSON
	var devices []pb.Device
	err = json.Unmarshal(data, &devices)
	if err != nil {
		return results, fmt.Errorf("erreur de décodage du JSON : %v", err)
	}

	// Parcour de chaque appareil
	for _, device := range devices {
		// Initialisation des structures pour compter les réussites et les échecs
		successCount := make(map[string]int32)
		failureCount := make(map[string]int32)

		for _, operation := range device.Operations {
			if operation.HasSucceeded {
				successCount[operation.Type]++
			} else {
				failureCount[operation.Type]++
			}
		}

		// Ajout des résultats pour cet appareil
		deviceResults := pb.DeviceResults{
			DeviceName:   device.DeviceName,
			SuccessCount: successCount,
			FailureCount: failureCount,
		}
		results = append(results, deviceResults)
	}

	return results, nil
}

/*
*	Fonction pour afficher les résultats
*	Résultat :
*	Device Name: 0b3939ec-06d4-48e0-ade9-d06e48fd4fe0
*	Success Counts:
*	  CREATE: 27
*	  DELETE: 13
*	  UPDATE: 20
*	Failure Counts:
*  	  DELETE: 6
*     UPDATE: 5
*     CREATE: 7
 */
func DisplayResults(results []pb.DeviceResults) {
	for _, deviceResults := range results {
		fmt.Printf("Device Name: %s\n", deviceResults.DeviceName)

		fmt.Println("Success Counts:")
		for opType, count := range deviceResults.SuccessCount {
			fmt.Printf("  %s: %d\n", opType, count)
		}

		fmt.Println("Failure Counts:")
		for opType, count := range deviceResults.FailureCount {
			fmt.Printf("  %s: %d\n", opType, count)
		}

		fmt.Println()
	}
}

/*== END/Public Functions  ==*/
