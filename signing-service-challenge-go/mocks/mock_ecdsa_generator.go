package mocks

import "signing-service-challenge/crypto"

type MockECCGenerator struct {
	GenerateFunc func() (*crypto.ECCKeyPair, error)
}

func (m *MockECCGenerator) Generate() (*crypto.ECCKeyPair, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc()
	}
	return nil, nil
}
