// package domain

// //import "github.com/google/uuid"
// import (
// 	base64 "encoding/base64"

// 	"fmt"

// 	uuid "github.com/google/uuid"
// )

// //signature device domain model ...
// type device struct {
// 	id                  uuid.UUID
// 	label               string
// 	signature_counter   int
// 	privateKeyBytes     []byte
// 	publicKeyBytes      []byte
// 	last_signature      string
// 	signature_algorithm SignatureAlgorithm
// }

// func New(id uuid.UUID, algorithm SignatureAlgorithm, label string) device {
// 	s := base64.RawURLEncoding.EncodeToString([]byte(id[:]))
// 	d := device{id, label, 0, nil, nil, s, algorithm}
// 	return d
// }

//	func (d device) Sign(data_to_sign string) string {
//		s := fmt.Sprint(d.signature_counter, data_to_sign, d.last_signature)
//		return s
//	}
package domain

import (
	"errors"

	"github.com/google/uuid"
)

type SignatureDevice struct {
	ID               uuid.UUID
	Label            string
	KeyPairAlgorithm CryptoAlgorithm
	SignatureCounter int
	PublicKey        string
	PrivateKey       string
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
	return &SignatureDevice{
		ID:               deviceid,
		Label:            label,
		KeyPairAlgorithm: CryptoAlgorithm,
		SignatureCounter: 0,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
	}, nil
}

func (s *SignatureDevice) IncrementSignatureCounter() {
	s.SignatureCounter++
}

// func (d *SignatureDevice) Sign(data_to_sign string) string {
// 	s := fmt.Sprint(d.SignatureCounter, data_to_sign, d.last_signature)
// 	return s
// }
