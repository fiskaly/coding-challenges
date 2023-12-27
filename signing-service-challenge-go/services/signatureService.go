package crypto

import (
	"slices"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

//TODO: Figure out how to make key marshallers work with interface

func CreateSignatureDevice(id string, algorithm string, label string) (domain.CreateSignatureDeviceResponse, error) {
	var privateKeyBytes []byte
	var publicKey string
	repo := persistence.New()

	//check if device exists
	_, err := repo.FindDeviceById(id)
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}, err
	}

	//chck that algoirthm is supported and generate keys
	algorithmUppercase := strings.ToUpper(algorithm)
	signatureAlgorithmRegistry := crypto.NewSignatureAlgorithmRegistry()
	if !slices.Contains(signatureAlgorithmRegistry.AlgorithmList, algorithm) {
		return domain.CreateSignatureDeviceResponse{}, err
	}
	if strings.Compare("RSA", algorithmUppercase) == 0 {
		privateKeyBytes, publicKey, err = generateRSAKeys()
	} else if strings.Compare("ECDSA", algorithmUppercase) == 0 {
		privateKeyBytes, publicKey, err = generateRSAKeys()
	}

	//create device and save it
	signatureDevice := domain.NewSignatureDevice(id, privateKeyBytes, publicKey, crypto.SignatureAlgorithm(algorithmUppercase), label)
	err = repo.NewDevice(*signatureDevice)
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}, err
	}

	return *signatureDevice.GetCreSignatureDeviceResponse(), nil
}

func SignTransaction(deviceId string, data string) {

}

func generateRSAKeys() ([]byte, string, error) {
	rsaGenerator := crypto.NewRSAGenerator()
	keypair, err := rsaGenerator.Generate()
	if err != nil {
		return nil, "", err
	}
	rsaMarshaler := crypto.NewRSAMarshaler()
	privateKeyBytes, _, err := rsaMarshaler.Marshal(*keypair)
	return privateKeyBytes, keypair.Public.N.String(), err
}

func generateECDSAKeys() ([]byte, []byte, error) {
	eccGenerator := crypto.NewECCGenerator()
	keypair, err := eccGenerator.Generate()
	if err != nil {
		return nil, nil, err
	}
	eccMarshaler := crypto.NewECCMarshaler()
	return eccMarshaler.Marshal(*keypair)
}
