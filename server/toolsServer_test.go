package main

import (
	"testing"

	pb "mahmoud.projet.rt0805/proto/SendData"
)

func TestVerifyData(t *testing.T) {
	// Cas de test avec données incorrectes (journée invalide et valeurs incorrectes)
	req := &pb.SendDataRequest{
		Journee: 0, // Journée invalide : doit être >= 1.
		DeviceResults: []*pb.DeviceResults{
			{
				DeviceName: "test-device",
				SuccessCount: map[string]int32{
					"CREATE":  5,
					"UPDATE":  -3, // Valeur négative : non autorisée.
					"INVALID": 10, // Clé non autorisée.
				},
				FailureCount: map[string]int32{
					"CREATE": -2, // Valeur négative : non autorisée.
					"UPDATE": 4,
				},
			},
		},
	}

	got := VerifyData(req)
	if got {
		t.Errorf("VerifyData() renvoie %v; attendu false pour les données incorrectes", got)
	}

	// Exemple de cas correct:
	req2 := &pb.SendDataRequest{
		Journee: 2,
		DeviceResults: []*pb.DeviceResults{
			{
				DeviceName: "test-device",
				SuccessCount: map[string]int32{
					"CREATE": 5,
					"UPDATE": 4,
					"DELETE": 3,
				},
				FailureCount: map[string]int32{
					"CREATE": 1,
					"UPDATE": 0,
					"DELETE": 2,
				},
			},
		},
	}

	got2 := VerifyData(req2)
	if !got2 {
		t.Errorf("VerifyData() renvoie %v; attendu true pour les données correctes", got2)
	}
}
