package crypto

import (
	"crypto/ecdsa"
	"testing"
)

func TestECCKeyPair_Sign(t *testing.T) {
	generator := ECCGenerator{}
	keyPair, err := generator.generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signature, err := keyPair.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	digest, err := computeHashDigest([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	result := ecdsa.VerifyASN1(keyPair.Public, digest, signature)
	if !result {
		t.Errorf("signature verification failed: %s", err)
	}
}
