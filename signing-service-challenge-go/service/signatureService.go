package service

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

//TODO: Figure out how to make key marshallers work with interface

func CreateSignatureDevice(device *domain.SignatureDevice) (domain.CreateSignatureDeviceResponse, error) {
	var privateKeyBytes, publicKey []byte
	repo := persistence.Get()

	//check if device exists
	_, err := repo.FindDeviceById(device.Id)
	if err == nil {
		return domain.CreateSignatureDeviceResponse{}, errors.New("[CreateSignatureDevice] device with specified ID already exists")
	}

	//chck that algoirthm is supported and generate keys
	signatureAlgorithmRegistry := crypto.NewSignatureAlgorithmRegistry()
	if !slices.Contains(signatureAlgorithmRegistry.AlgorithmList, device.Algorithm) {
		return domain.CreateSignatureDeviceResponse{}, errors.New("algorithm, not supported: " + device.Algorithm)
	}
	if strings.Compare(signatureAlgorithmRegistry.RSA, device.Algorithm) == 0 {
		privateKeyBytes, publicKey, err = generateRSAKeys()
	} else if strings.Compare(signatureAlgorithmRegistry.ECDSA, device.Algorithm) == 0 {
		privateKeyBytes, publicKey, err = generateECDSAKeys()
	}
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}, err
	}

	//create device and save it
	device.PrivateKeyBytes = privateKeyBytes
	device.PublicKey = publicKey
	err = repo.NewDevice(*device)
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}, err
	}

	return *device.GetCreateSignatureDeviceResponse(), nil
}

func SignTransaction(deviceId string, data string) (domain.CreateSignatureResponse, error) {
	repo := persistence.Get()
	signatureDevice, err := repo.FindDeviceById(deviceId)
	if err != nil {
		return domain.CreateSignatureResponse{}, err
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
	signatureAlgorithmRegistry := crypto.NewSignatureAlgorithmRegistry()
	var signer crypto.Signer
	if strings.Compare(signatureAlgorithmRegistry.RSA, signatureDevice.Algorithm) == 0 {
		signer = crypto.NewRSASigner(signatureDevice.PrivateKeyBytes)
	} else if strings.Compare(signatureAlgorithmRegistry.ECDSA, signatureDevice.Algorithm) == 0 {
		signer = crypto.NewECCSigner(signatureDevice.PrivateKeyBytes)
	} else {
		return domain.CreateSignatureResponse{}, errors.New("algorithm, not supported: " + signatureDevice.Algorithm)
	}
	signature, err := signer.Sign([]byte(signing_string))
	if err != nil {
		return domain.CreateSignatureResponse{}, err
	}

	//update device and save it
	signatureDevice.LastSignature = signature
	signatureDevice.SignatureCounter += 1
	err = repo.UpdateDevice(signatureDevice)
	if err != nil {
		return domain.CreateSignatureResponse{}, err
	}
	return *signatureDevice.GetSignatureResponse(), nil

}

func GetDeviceInfo(deviceId string) (domain.SignatureDeviceInfoResponse, error) {
	repo := persistence.Get()
	device, err := repo.FindDeviceById(deviceId)
	if err != nil {
		return domain.SignatureDeviceInfoResponse{}, err
	}
	return *device.GetSignatureDeviceInfoResponse(), err
}

func generateRSAKeys() ([]byte, []byte, error) {
	rsaGenerator := crypto.NewRSAGenerator()
	keypair, err := rsaGenerator.Generate()
	if err != nil {
		return nil, nil, err
	}
	rsaMarshaler := crypto.NewRSAMarshaler()
	_, privateKeyBytes, err := rsaMarshaler.Marshal(*keypair)

	publicKeyBytes := x509.MarshalPKCS1PublicKey(keypair.Public)
	return privateKeyBytes, publicKeyBytes, err
}

func generateECDSAKeys() ([]byte, []byte, error) {
	eccGenerator := crypto.NewECCGenerator()
	keypair, err := eccGenerator.Generate()
	if err != nil {
		return nil, nil, err
	}
	eccMarshaler := crypto.NewECCMarshaler()
	_, privateKeyBytes, err := eccMarshaler.Marshal(*keypair)
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(keypair.Public)
	return privateKeyBytes, publicKeyBytes, err
}
