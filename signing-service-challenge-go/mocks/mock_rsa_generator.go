package mocks

import "signing-service-challenge/crypto"

type MockRSAGenerator struct {
	GenerateFunc func() (*crypto.RSAKeyPair, error)
}

func (m *MockRSAGenerator) Generate() (*crypto.RSAKeyPair, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc()
	}
	return nil, nil
}
