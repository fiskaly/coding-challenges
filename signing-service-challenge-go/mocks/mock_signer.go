package mocks

type MockSigner struct {
	SignFunc func(dataToBeSigned []byte) ([]byte, error)
}

func (m *MockSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	if m.SignFunc != nil {
		return m.SignFunc(dataToBeSigned)
	}
	return nil, nil
}
