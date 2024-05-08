package main

import (
	pb "mahmoud.projet.rt0805/proto/SendData"
)

/*
* VerifyData reçoit une requête SendDataRequest et vérifie que
* chaque DeviceResult contient les clés nécessaires
* (CREATE, UPDATE, DELETE) dans SuccessCount et FailureCount.
 */
func VerifyData(req *pb.SendDataRequest) bool {
	requiredKeys := []string{"CREATE", "UPDATE", "DELETE"}

	if req.Journee < 1 {
		return false
	}

	for _, deviceResult := range req.DeviceResults {
		if deviceResult == nil {
			return false
		}

		missingKeys := 0
		for _, key := range requiredKeys {
			successValue, successExists := deviceResult.SuccessCount[key]
			if !successExists {
				missingKeys++
			} else {
				if successValue < 0 {
					return false
				}
			}

			failureValue, failureExists := deviceResult.FailureCount[key]
			if failureExists && failureValue < 0 {
				return false
			}
		}

		if missingKeys > 2 {
			return false
		}

		for key := range deviceResult.SuccessCount {
			if key != "CREATE" && key != "UPDATE" && key != "DELETE" {
				return false
			}
		}

		for key := range deviceResult.FailureCount {
			if key != "CREATE" && key != "UPDATE" && key != "DELETE" {
				return false
			}
		}
	}

	return true
}

// func VerifyData(req *pb.SendDataRequest) bool {
// 	// Liste des clés nécessaires dans SuccessCount et FailureCount.
// 	requiredKeys := []string{"CREATE", "UPDATE", "DELETE"}

// 	// Vérifier que la journée est supérieure ou égale à 1.
// 	if req.Journee < 1 {
// 		log.Println("La journée doit être supérieure ou égale à 1.")
// 		return false
// 	}

// 	// Parcourt les DeviceResults dans la requête.
// 	for _, deviceResult := range req.DeviceResults {
// 		// Vérifier que deviceResult n'est pas nil.
// 		if deviceResult == nil {
// 			log.Println("Un deviceResult est nil.")
// 			return false
// 		}

// 		// Compte le nombre de clés manquantes dans SuccessCount.
// 		missingKeys := 0

// 		// Vérifier que SuccessCount contient au maximum 2 clés manquantes.
// 		for _, key := range requiredKeys {
// 			// Vérifier SuccessCount
// 			successValue, successExists := deviceResult.SuccessCount[key]
// 			if !successExists {
// 				missingKeys++
// 			} else {
// 				if successValue < 0 {
// 					log.Printf("La valeur de %s dans SuccessCount est négative : %d", key, successValue)
// 					return false
// 				}
// 			}

// 			// Vérifier FailureCount
// 			failureValue, failureExists := deviceResult.FailureCount[key]
// 			if failureExists && failureValue < 0 {
// 				log.Printf("La valeur de %s dans FailureCount est négative : %d", key, failureValue)
// 				return false
// 			}
// 		}

// 		// Vérifier le nombre de clés manquantes dans SuccessCount.
// 		if missingKeys > 2 {
// 			log.Printf("Plus de 2 clés manquantes dans SuccessCount : %d", missingKeys)
// 			return false
// 		}

// 		// Vérifier qu'il n'y a pas de clés non autorisées dans SuccessCount et FailureCount.
// 		for key := range deviceResult.SuccessCount {
// 			if key != "CREATE" && key != "UPDATE" && key != "DELETE" {
// 				log.Printf("Clé non autorisée dans SuccessCount : %s", key)
// 				return false
// 			}
// 		}

// 		for key := range deviceResult.FailureCount {
// 			if key != "CREATE" && key != "UPDATE" && key != "DELETE" {
// 				log.Printf("Clé non autorisée dans FailureCount : %s", key)
// 				return false
// 			}
// 		}
// 	}

// 	// Si tous les DeviceResults respectent les conditions, retourne true.
// 	return true
// }
