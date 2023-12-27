package crypto

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

//TODO: Figure out how to make key marshallers work with interface

func CreateSignatureDevice(id string, algorithm string, label string) (api.CreateSignatureDeviceResponse, error) {
	var privateKeyBytes []byte
	var publicKey string
	repo := persistence.New()

	//check if device exists
	_, err := repo.FindDeviceById(id)
	if err != nil {
		return api.CreateSignatureDeviceResponse{}, err
	}

	//chck that algoirthm is supported and generate keys
	algorithmUppercase := strings.ToUpper(algorithm)
	signatureAlgorithmRegistry := crypto.NewSignatureAlgorithmRegistry()
	if !slices.Contains(signatureAlgorithmRegistry.AlgorithmList, algorithm) {
		return api.CreateSignatureDeviceResponse{}, err
	}
	if strings.Compare("RSA", algorithmUppercase) == 0 {
		privateKeyBytes, publicKey, err = generateRSAKeys()
	} else if strings.Compare("ECDSA", algorithmUppercase) == 0 {
		privateKeyBytes, publicKey, err = generateECDSAKeys()
	}

	//create device and save it
	signatureDevice := domain.NewSignatureDevice(id, privateKeyBytes, publicKey, crypto.SignatureAlgorithm(algorithmUppercase), label)
	err = repo.NewDevice(*signatureDevice)
	if err != nil {
		return api.CreateSignatureDeviceResponse{}, err
	}

	return *signatureDevice.GetCreSignatureDeviceResponse(), nil
}

func SignTransaction(deviceId string, data string) (api.SignatureResponse, error) {
	repo := persistence.New()
	signatureDevice, err := repo.FindDeviceById(deviceId)
	if err != nil {
		return api.SignatureResponse{}, err
	}
	//build signing string
	var part1, part2, part3 string
	part1 = fmt.Sprint(signatureDevice.SignatureCounter)
	part2 = data
	if signatureDevice.SignatureCounter == 0 {
		part3 = base64.StdEncoding.EncodeToString([]byte(deviceId))
	} else {
		part3 = string(signatureDevice.LastSignature)
	}
	signing_string := fmt.Sprintf("%s_%s_%s", part1, part2, part3)

	//get signer and sign
	var signer crypto.Signer
	if strings.Compare("RSA", string(signatureDevice.Algorithm)) == 0 {
		signer = crypto.NewRSASigner(signatureDevice.PrivateKeyBytes)
	} else if strings.Compare("ECDSA", string(signatureDevice.Algorithm)) == 0 {
		signer = crypto.NewECCSigner(signatureDevice.PrivateKeyBytes)
	}
	signature, err := signer.Sign([]byte(signing_string))
	if err != nil {
		return api.SignatureResponse{}, err
	}

	//update device and save it
	signatureDevice.LastSignature = signature
	err = repo.UpdateDevice(signatureDevice)
	if err != nil {
		return api.SignatureResponse{}, err
	}
	return *signatureDevice.GetSignatureResponse(), nil

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

func generateECDSAKeys() ([]byte, string, error) {
	eccGenerator := crypto.NewECCGenerator()
	keypair, err := eccGenerator.Generate()
	if err != nil {
		return nil, "", err
	}
	eccMarshaler := crypto.NewECCMarshaler()
	privateKeyBytes, _, err := eccMarshaler.Marshal(*keypair)
	return privateKeyBytes, keypair.Public.X.String(), err
}
