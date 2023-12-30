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

type SignatureServiceImplementation struct {
	rsaGenerator   crypto.RSAGenerator
	ecdsaGenerator crypto.ECCGenerator
	repo           persistence.DeviceRepository
}

func New() *SignatureServiceImplementation {
	return &SignatureServiceImplementation{
		rsaGenerator:   crypto.NewRSAGenerator(),
		ecdsaGenerator: crypto.NewECCGenerator(),
		repo:           persistence.Get(),
	}
}

// CreateSignatureDevice will create a new domain.SignatureDevice based on the supplied signature device.
// The keys for this device will be generated based on the selected signing algorithm.
// Returns a domain.CreateSignatureDeviceResponse entity ready to be sent in API response
func (service *SignatureServiceImplementation) CreateSignatureDevice(device *domain.SignatureDevice) (domain.CreateSignatureDeviceResponse, error) {
	var privateKeyBytes, publicKey []byte
	var err error

	//check if input is valid
	if len(device.Id) < 3 {
		return domain.CreateSignatureDeviceResponse{}, errors.New("[CreateSignatureDevice] device ID must be at least 3 characters long")
	}
	if len(device.Algorithm) < 1 {
		return domain.CreateSignatureDeviceResponse{}, errors.New("[CreateSignatureDevice] please specify a cryptographic algrithm for the device")
	}

	//check if device exists
	_, exists := service.repo.FindDeviceById(device.Id)
	if exists {
		return domain.CreateSignatureDeviceResponse{}, errors.New("[CreateSignatureDevice] device with specified ID already exists")
	}

	//chck that algoirthm is supported and generate keys
	signatureAlgorithmRegistry := crypto.NewSignatureAlgorithmRegistry()
	if !slices.Contains(signatureAlgorithmRegistry.AlgorithmList, device.Algorithm) {
		return domain.CreateSignatureDeviceResponse{}, errors.New("algorithm, not supported: " + device.Algorithm)
	}
	if strings.Compare(signatureAlgorithmRegistry.RSA, device.Algorithm) == 0 {
		privateKeyBytes, publicKey, err = service.generateRSAKeys()
	} else if strings.Compare(signatureAlgorithmRegistry.ECDSA, device.Algorithm) == 0 {
		privateKeyBytes, publicKey, err = service.generateECDSAKeys()
	}
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}, err
	}

	//create device and save it
	device.PrivateKeyBytes = privateKeyBytes
	device.PublicKey = publicKey
	err = service.repo.NewDevice(*device)
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}, err
	}

	return *device.GetCreateSignatureDeviceResponse(), nil
}

// generates an RSA keypair
// returns private key, public key
func (service *SignatureServiceImplementation) generateRSAKeys() ([]byte, []byte, error) {
	keypair, err := service.rsaGenerator.Generate()
	if err != nil {
		return nil, nil, err
	}
	rsaMarshaler := crypto.NewRSAMarshaler()
	_, privateKeyBytes, err := rsaMarshaler.Marshal(*keypair)

	publicKeyBytes := x509.MarshalPKCS1PublicKey(keypair.Public)
	return privateKeyBytes, publicKeyBytes, err
}

// generates an ECDSA keypair
// returns private key, public key
func (service *SignatureServiceImplementation) generateECDSAKeys() ([]byte, []byte, error) {
	keypair, err := service.ecdsaGenerator.Generate()
	if err != nil {
		return nil, nil, err
	}
	eccMarshaler := crypto.NewECCMarshaler()
	_, privateKeyBytes, err := eccMarshaler.Marshal(*keypair)
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(keypair.Public)
	return privateKeyBytes, publicKeyBytes, err
}

// Signs the supplied data with the device.
// Returns a domain.CreateSignatureResponse entity ready to be sent in API response
func (service *SignatureServiceImplementation) SignTransaction(deviceId string, data string) (domain.CreateSignatureResponse, error) {
	signatureDevice, exists := service.repo.FindDeviceById(deviceId)
	if !exists {
		return domain.CreateSignatureResponse{}, fmt.Errorf("[FindDeviceById] device with specified ID doesn't exist: \"%s\"", deviceId)
	}
	//build signing string
	signing_string := service.buildSigningString(signatureDevice, data)

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
	err = service.repo.UpdateDevice(signatureDevice)
	if err != nil {
		return domain.CreateSignatureResponse{}, err
	}
	return *signatureDevice.GetSignatureResponse(signing_string), nil
}

// Builds the signing string for a device & data combo
func (service *SignatureServiceImplementation) buildSigningString(signatureDevice domain.SignatureDevice, data string) string {
	var part1, part2, part3 string
	part1 = fmt.Sprint(signatureDevice.SignatureCounter)
	part2 = data
	if signatureDevice.SignatureCounter == 0 {
		part3 = base64.StdEncoding.EncodeToString([]byte(signatureDevice.Id))
	} else {
		part3 = base64.StdEncoding.EncodeToString(signatureDevice.LastSignature)
	}
	return fmt.Sprintf("%s_%s_%s", part1, part2, part3)
}

// Returns all the information on a specified device in the form of domain.SignatureDeviceInfoResponse, ready to be returned by API
func (service *SignatureServiceImplementation) GetDeviceInfo(deviceId string) (domain.SignatureDeviceInfoResponse, error) {
	device, exists := service.repo.FindDeviceById(deviceId)
	if !exists {
		return domain.SignatureDeviceInfoResponse{}, fmt.Errorf("[FindDeviceById] device with specified ID doesn't exist: \"%s\"", deviceId)
	}
	return *device.GetSignatureDeviceInfoResponse(), nil
}

func (service *SignatureServiceImplementation) GetAllDevices() []string {
	return service.repo.GetAllDevices()
}
