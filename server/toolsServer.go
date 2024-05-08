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
