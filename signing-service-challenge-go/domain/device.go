package domain

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/util"
)

type SignatureDevice struct {
	Id               string
	Label            string
	Algorithm        string
	signer           crypto.Signer
	signatureCounter uint
	lastSignature    string
	mutex            sync.Mutex
}

type SignDataResult struct {
	Signature  []byte
	SignedData []byte
}

func NewSignatureDevice(
	id string,
	label string,
	algorithm string,
	signer crypto.Signer,
) *SignatureDevice {
	lastSignature := util.EncodeToBase64String([]byte(id))
	return &SignatureDevice{
		Id:               id,
		Label:            label,
		Algorithm:        algorithm,
		signer:           signer,
		signatureCounter: 0,
		lastSignature:    lastSignature,
	}
}

func (device *SignatureDevice) Sign(dataToBeSigned []byte) (*SignDataResult, error) {
	device.mutex.Lock()
	defer device.mutex.Unlock()

	securedData := device.secureData(dataToBeSigned)

	signature, err := device.signer.Sign(securedData)
	if err != nil {
		return nil, err
	}

	device.signatureCounter += 1
	device.lastSignature = util.EncodeToBase64String(signature)

	return &SignDataResult{
		Signature:  signature,
		SignedData: securedData,
	}, nil
}

func (device *SignatureDevice) secureData(dataToBeSigned []byte) []byte {
	sigCounter := []byte(fmt.Sprintf("%d_", device.signatureCounter))
	lastSig := []byte(fmt.Sprintf("_%s", device.lastSignature))
	securedData := append(append(sigCounter, dataToBeSigned...), lastSig...)

	return securedData
}
