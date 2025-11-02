package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"testing"
)

func TestRSASignerSignsData(t *testing.T) {
	generator := &RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("failed to generate RSA key pair: %v", err)
	}

	signer, err := NewRSASigner(keyPair)
	if err != nil {
		t.Fatalf("expected no error creating signer, got %v", err)
	}

	data := []byte("payload-to-sign")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("expected no error signing data, got %v", err)
	}

	hashed := sha256.Sum256(data)
	if err := rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hashed[:], signature); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestECDSASignerSignsData(t *testing.T) {
	generator := &ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("failed to generate ECC key pair: %v", err)
	}

	signer, err := NewECDSASigner(keyPair)
	if err != nil {
		t.Fatalf("expected no error creating ECDSA signer, got %v", err)
	}

	data := []byte("payload-to-sign")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("expected no error signing data, got %v", err)
	}

	hasher := sha512.New384()
	if _, err := hasher.Write(data); err != nil {
		t.Fatalf("failed to hash data: %v", err)
	}

	digest := hasher.Sum(nil)
	if !ecdsa.VerifyASN1(keyPair.Public, digest, signature) {
		t.Fatal("signature verification failed")
	}
}

func TestNewSignersValidateKeyPairs(t *testing.T) {
	if _, err := NewRSASigner(nil); !errors.Is(err, ErrNilKeyPair) {
		t.Fatalf("expected ErrNilKeyPair, got %v", err)
	}

	if _, err := NewECDSASigner(nil); !errors.Is(err, ErrNilKeyPair) {
		t.Fatalf("expected ErrNilKeyPair, got %v", err)
	}

	if _, err := NewRSASigner(&RSAKeyPair{}); !errors.Is(err, ErrNilPrivateKey) {
		t.Fatalf("expected ErrNilPrivateKey, got %v", err)
	}

	if _, err := NewECDSASigner(&ECCKeyPair{}); !errors.Is(err, ErrNilPrivateKey) {
		t.Fatalf("expected ErrNilPrivateKey, got %v", err)
	}
}

func TestSignErrorsWhenSignerUninitialised(t *testing.T) {
	signer := &RSASigner{}
	if _, err := signer.Sign([]byte("data")); !errors.Is(err, ErrNilPrivateKey) {
		t.Fatalf("expected ErrNilPrivateKey, got %v", err)
	}

	ecdsaSigner := &ECDSASigner{}
	if _, err := ecdsaSigner.Sign([]byte("data")); !errors.Is(err, ErrNilPrivateKey) {
		t.Fatalf("expected ErrNilPrivateKey, got %v", err)
	}
}
