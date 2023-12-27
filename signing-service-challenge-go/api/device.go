package api

import "github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"

// TODO: REST endpoints ...
type CreateSignatureDeviceResponse struct {
	DeviceId  string
	PublicKey string
	Algorithm crypto.SignatureAlgorithm
	Alias     string
}

type SignatureResponse struct {
	DeviceId          string
	Signature         string
	SignaturesCreated int64
	PublicKey         string
	Algorithm         crypto.SignatureAlgorithm
	Alias             string
}

type SignatureDeviceInfoResponse struct {
	DeviceId          string
	PublicKey         string
	Algorithm         crypto.SignatureAlgorithm
	SignaturesCreated int64
	LastSignature     string
	Alias             string
}
