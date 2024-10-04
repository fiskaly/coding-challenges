package domain

import (
	"encoding/base64"
	"fmt"
	"signing-service-challenge/crypto"
	"sync"
)

type AlgorithmType string

const (
	RSAAlgorithm AlgorithmType = "RSA"
	ECCAlgorithm AlgorithmType = "ECC"
)

// TODO: signature device domain model ...
type Device struct {
	mu               sync.RWMutex
	Id               string
	Label            string
	SignatureCounter int
	Algorithm        AlgorithmType
	Signer           crypto.Signer
	LastSignature    []byte
	PublicKey        []byte
	PrivateKey       []byte
}

func (device *Device) GetSignatureReference() []byte {
	device.mu.RLock()
	defer device.mu.RUnlock()
	if device.SignatureCounter == 0 {
		return []byte(device.Id)
	}
	return device.LastSignature
}

func (device *Device) BuildSecuredDataToBeSigned(data string) string {
	signatureReference := device.GetSignatureReference()
	encodedSignature := base64.StdEncoding.EncodeToString(signatureReference)
	device.mu.RLock()
	securedDataToBeSigned := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, encodedSignature)
	device.mu.RUnlock()
	return securedDataToBeSigned
}

func (device *Device) IncrementSignatureCounter() {
	device.mu.Lock()
	defer device.mu.Unlock()
	device.SignatureCounter++
}

func (devide *Device) UpdateLastSignature(signature []byte) {
	devide.mu.Lock()
	defer devide.mu.Unlock()
	devide.LastSignature = signature
}
