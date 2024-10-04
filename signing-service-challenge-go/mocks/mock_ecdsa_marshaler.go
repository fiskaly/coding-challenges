package mocks

import "signing-service-challenge/crypto"

type MockECCMarshaler struct {
	EncodeFunc func(keyPair crypto.ECCKeyPair) ([]byte, []byte, error)
	DecodeFunc func(privateKeyBytes []byte) (*crypto.ECCKeyPair, error)
}

func (m *MockECCMarshaler) Encode(keyPair crypto.ECCKeyPair) ([]byte, []byte, error) {
	if m.EncodeFunc != nil {
		return m.EncodeFunc(keyPair)
	}
	return nil, nil, nil
}

func (m *MockECCMarshaler) Decode(privateKeyBytes []byte) (*crypto.ECCKeyPair, error) {
	if m.DecodeFunc != nil {
		return m.DecodeFunc(privateKeyBytes)
	}
	return nil, nil
}
