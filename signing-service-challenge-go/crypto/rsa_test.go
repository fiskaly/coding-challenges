package crypto

import "testing"

func TestRSAAlgorithm_Name(t *testing.T) {
	rsa := RSAAlgorithm{}
	name := rsa.Name()

	if name != "RSA" {
		t.Errorf("expected RSA, got: %s", name)
	}
}

func TestRSAAlgorithm_GenerateEncodedPrivateKey(t *testing.T) {
	rsa := RSAAlgorithm{}
	encodedPrivateKey, err := rsa.GenerateEncodedPrivateKey()

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	_, err = NewRSAMarshaler().Unmarshal(encodedPrivateKey)
	if err != nil {
		t.Errorf("decode of generated private key failed: %s", err)
	}
}
