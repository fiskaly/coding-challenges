package domain

import (
	"encoding/base64"
	"strconv"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type SignatureDevice struct {
	id               string
	label            string
	algorithm        string
	publicKey        any
	privateKey       any
	signatureCounter int32
	lastSignature    string
	mu               sync.Mutex
}

func CreateNewDevice(algorithm, label string) (*SignatureDevice, error) {
	adapter, err := crypto.NewGeneratorAdapter(algorithm)
	if err != nil {
		return nil, err
	}
	keyPair, err := adapter.Generate()
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	return &SignatureDevice{
		id:               id,
		label:            label,
		algorithm:        algorithm,
		publicKey:        keyPair.PublicKey,
		privateKey:       keyPair.PrivateKey,
		signatureCounter: 0,
		lastSignature:    base64.StdEncoding.EncodeToString([]byte(id)),
	}, nil
}

func (d *SignatureDevice) SignData(data string) (string, string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	securedData := d.GetSecuredDataToBeSigned(data)
	signer, err := crypto.NewSignerWithKey(d.algorithm, d.privateKey)
	if err != nil {
		return "", "", err
	}

	signature, err := signer.Sign([]byte(securedData))
	if err != nil {
		return "", "", err
	}

	d.lastSignature = signature
	d.signatureCounter++
	return signature, securedData, nil
}

func (d *SignatureDevice) GetSecuredDataToBeSigned(dataToBeSigned string) string {
	counterStr := strconv.Itoa(int(d.signatureCounter))
	return counterStr + "_" + dataToBeSigned + "_" + d.lastSignature
}

func (d *SignatureDevice) GetSignatureCounter() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return int(d.signatureCounter)
}

func (d *SignatureDevice) GetID() string {
	return d.id
}

func (d *SignatureDevice) GetAlgorithm() string {
	return d.algorithm
}

func (d *SignatureDevice) GetLabel() string {
	return d.label
}
