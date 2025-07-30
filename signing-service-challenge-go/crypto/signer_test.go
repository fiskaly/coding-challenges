package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestRSASigner(t *testing.T) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	signer := &RSASigner{privateKey: *privateKey}
	data := []byte("test data")

	// Sign the data
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("failed to sign data: %v", err)
	}

	// Verify the signature
	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(&privateKey.PublicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		t.Fatalf("failed to verify signature: %v", err)
	}
}

func TestECCSigner(t *testing.T) {
	// Generate ECC key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ECC key: %v", err)
	}

	signer := &ECCSigner{privateKey: *privateKey}
	data := []byte("test data")

	// Sign the data
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("failed to sign data: %v", err)
	}

	// Verify the signature
	hashed := sha256.Sum256(data)
	if !ecdsa.VerifyASN1(&privateKey.PublicKey, hashed[:], signature) {
		t.Fatalf("failed to verify signature")
	}
}
