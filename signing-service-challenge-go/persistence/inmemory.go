package persistence

// TODO: in-memory persistence ...
// SQL can also be used I did not do it cause I was running out of time. My Initial idea was SQL database.
import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func SaveDevice(device domain.SignatureDevice, filePath string) error {
	jsonData, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SignatureDevice saved to %s\n", filePath)
	return nil
}
