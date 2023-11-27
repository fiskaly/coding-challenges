package domain

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type EncodedKeyPair struct {
}

type SignatureDevice struct {
	Id               string
	Label            *string
	signer           crypto.Signer
	signatureCounter uint
	lastSignature    string
}

type SignDataResult struct {
	Signature  []byte
	SignedData []byte
}

func NewSignatureDevice(id string, label *string, algorithm string) (*SignatureDevice, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid id")
	}

	// 1. create generator based on algorithm
	// 2. create marshaller based on algorithm
	// 3. marshal key pair
	// 4. store marshaller and key pair in signer
	// 5. store signer in signature device

	signer, err := crypto.NewSigner(algorithm)
	if err != nil {
		return nil, err
	}

	// 1. check that device with this id already exists
	// 2. if yes, return error
	// 3. if no, persist

	lastSignature := base64.StdEncoding.EncodeToString([]byte(id))

	return &SignatureDevice{
		Id:               id,
		Label:            label,
		signer:           signer,
		signatureCounter: 0,
		lastSignature:    lastSignature,
	}, nil
}

func (device *SignatureDevice) Sign(dataToBeSigned []byte) (*SignDataResult, error) {
	securedData := device.secureData(dataToBeSigned)

	signature, err := device.signer.Sign(securedData)
	if err != nil {
		return nil, err
	}

	// TODO: make thread safe
	device.signatureCounter += 1
	device.lastSignature = base64.StdEncoding.EncodeToString(signature)

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
