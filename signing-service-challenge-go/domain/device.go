package domain

import (
	"crypto"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto/generator"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto/signer"
	"github.com/google/uuid"
)

// -- Device
// unique identifier	-> e.g uuid
// signature algorithm 	-> that the device will be using to sign transaction data
// key pair 			-> during the creation process, new heypair has to be generated and assigned to the device
// label				-> used to display a label in the ui
// signatureCounter		-> tracks how many signature have been created

// Created_at			-> optional (date of the device creation)

// ---------------------------------------------------------------------------

// -- Signature creation

// client will have to provide the data_to_be_signed through the API
// 		- to increase the security of the system we will extend this raw data with
//		signature_counter and last signature
//
//		The resulting string (secured_data_to_be_signed) should follow this format:
//		 <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
//
// 		In the base case there is no last_signature (= signature_counter == 0).
//		 Use the base64-encoded device ID (last_signature = base64(device.id)) instead of the last_signature.

// ---------------------------------------------------------------------------
// CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse
// SignTransaction(deviceId: string, data: string): SignatureResponse

// SignatureDevice represents of generating keys and signing data.
type SignatureDevice struct {
	ID               string
	Algorithm        string
	Label            string
	KeyPair          crypto.PrivateKey
	SignatureCounter int64
	CreatedAt        time.Time
}

func NewSignatureDevice(id string, algorithm string, label string) (*SignatureDevice, error) {
	gen, err := generator.GetGenerator(algorithm)
	if err != nil {
		return nil, err
	}

	keyPair, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	return &SignatureDevice{
		ID:               id,
		Algorithm:        algorithm,
		Label:            label,
		KeyPair:          keyPair,
		SignatureCounter: 0,
		CreatedAt:        time.Now(),
	}, nil
}

func (sg *SignatureDevice) SignTransaction(data string) (*Transaction, error) {

	signer, err := signer.GetSigner(sg.Algorithm)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign([]byte(data))
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:        uuid.New().String(),
		DeviceID:  sg.ID,
		Data:      data,
		Signature: signature,
		Timestamp: time.Now(),
	}, nil
}
