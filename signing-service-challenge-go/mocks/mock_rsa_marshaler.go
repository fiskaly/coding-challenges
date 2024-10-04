package mocks

import "signing-service-challenge/crypto"

type MockRSAMarshaler struct {
	MarshalFunc   func(keyPair crypto.RSAKeyPair) ([]byte, []byte, error)
	UnmarshalFunc func(privateKeyBytes []byte) (*crypto.RSAKeyPair, error)
}

func (m *MockRSAMarshaler) Marshal(keyPair crypto.RSAKeyPair) ([]byte, []byte, error) {
	if m.MarshalFunc != nil {
		return m.MarshalFunc(keyPair)
	}
	return nil, nil, nil
}

func (m *MockRSAMarshaler) Unmarshal(privateKeyBytes []byte) (*crypto.RSAKeyPair, error) {
	if m.UnmarshalFunc != nil {
		return m.UnmarshalFunc(privateKeyBytes)
	}
	return nil, nil
}
