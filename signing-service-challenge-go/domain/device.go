package domain

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SignatureDevice struct {
	ID               uuid.UUID
	Label            string
	KeyPairAlgorithm CryptoAlgorithm
	signatureCounter int
	PublicKey        string
	privateKey       string
	last_signature   string
}

func NewSignatureDevice(id, label, algorithm, publicKey, privateKey string) (*SignatureDevice, error) {

	var deviceid uuid.UUID
	if uid, err := uuid.Parse(id); err == nil {
		deviceid = uid
	} else {
		return nil, errors.New("invalid device ID")
	}

	// Validate algorithm
	CryptoAlgorithm, err := FromString(algorithm)
	if err != nil {
		return nil, errors.New("invalid algorithm")
	}
	lastsignature := base64Encode(deviceid.String())
	return &SignatureDevice{
		ID:               deviceid,
		Label:            label,
		KeyPairAlgorithm: CryptoAlgorithm,
		signatureCounter: 0,
		PublicKey:        publicKey,
		privateKey:       privateKey,
		last_signature:   lastsignature,
	}, nil
}

func (s *SignatureDevice) incrementSignatureCounter() {
	s.signatureCounter++
}
func (s *SignatureDevice) updateLastSignature(lastsign string) {
	s.last_signature = lastsign
}
func (s *SignatureDevice) formatDataToSign(data_to_sign string) (signeddata string) {

	sdata := fmt.Sprint(s.signatureCounter, data_to_sign, s.last_signature)
	return sdata
}
func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
func (d *SignatureDevice) Sign(data_to_sign string, signer Signer) (*SignTransaction, error) {
	data_to_sign = d.formatDataToSign(data_to_sign)
	signeddata, err := signer.Sign(d.privateKey, data_to_sign)
	if err != nil {
		return nil, errors.New("signature failure")
	}
	d.incrementSignatureCounter()
	signenc := base64Encode(signeddata)
	d.updateLastSignature(signenc)
	sign, err := NewSignTransaction(d.ID.String(), data_to_sign, signenc)

	return sign, err
}
