package domain

import (
	"fmt"

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
}

type SignDataResult struct {
	Signature  []byte
	SignedData []byte
}

func (device *SignatureDevice) Sign(dataToBeSigned []byte) (*SignDataResult, error) {
	securedData := device.secureData(dataToBeSigned)

	signature, err := device.signer.Sign(securedData)
	if err != nil {
		return nil, err
	}

	// TODO: make thread safe
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
