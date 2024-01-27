package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"testing"
)

func TestRSASigner_Sign(t *testing.T) {
	generator := RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signer := RSASigner{keyPair: *keyPair}
	signature, err := signer.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	digest, err := computeDigestWithHashFunction([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	err = rsa.VerifyPSS(keyPair.Public, hashFunction, digest, signature, nil)
	if err != nil {
		t.Errorf("signature verification failed: %s", err)
	}
}

func TestECCSigner_Sign(t *testing.T) {
	generator := ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatal(err)
	}

	dataToBeSigned := "some-data"
	signer := ECCSigner{keyPair: *keyPair}
	signature, err := signer.Sign([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	digest, err := computeDigestWithHashFunction([]byte(dataToBeSigned))
	if err != nil {
		t.Fatal(err)
	}

	result := ecdsa.VerifyASN1(keyPair.Public, digest, signature)
	if !result {
		t.Errorf("signature verification failed: %s", err)
	}
}
