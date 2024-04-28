package main

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"mahmoud.projet.rt0805/proto/SendData"
)

func TestParseFile(t *testing.T) {
	// Crée un répertoire temporaire
	tempDir := t.TempDir()

	// Chemin complet vers le fichier de test JSON dans le répertoire temporaire
	testFilePath := filepath.Join(tempDir, "test.json")

	// Données JSON pour le test
	testData := `[{"device_name":"device_1","operations":[{"type":"CREATE","has_succeeded":true},{"type":"DELETE","has_succeeded":false},{"type":"UPDATE","has_succeeded":true}]},{"device_name":"device_2","operations":[{"type":"UPDATE","has_succeeded":true},{"type":"CREATE","has_succeeded":false},{"type":"DELETE","has_succeeded":true}]}]`

	// Écrit les données JSON dans le fichier de test
	err := ioutil.WriteFile(testFilePath, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier de test JSON : %v", err)
	}

	// Appel de ParseFile avec le chemin complet du fichier de test JSON
	results, err := ParseFile(testFilePath)
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution de ParseFile : %v", err)
	}

	// Résultats attendus
	expectedResults := []SendData.DeviceResults{
		{
			DeviceName: "device_1",
			SuccessCount: map[string]int32{
				"CREATE": 1,
				"UPDATE": 1,
			},
			FailureCount: map[string]int32{
				"DELETE": 1,
			},
		},
		{
			DeviceName: "device_2",
			SuccessCount: map[string]int32{
				"UPDATE": 1,
				"DELETE": 1,
			},
			FailureCount: map[string]int32{
				"CREATE": 1,
			},
		},
	}

	// Compare les résultats obtenus aux résultats attendus
	for i, result := range results {
		if !reflect.DeepEqual(&result, &expectedResults[i]) {
			t.Errorf("Les résultats obtenus ne correspondent pas aux résultats attendus. Obtenu: %v, Attendu: %v", result, expectedResults[i])
		}
	}
}
