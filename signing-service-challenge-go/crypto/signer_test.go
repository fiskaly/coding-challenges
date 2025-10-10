package crypto

import (
	"testing"
)

// TestRSASigner tests RSA signing and verification
func TestRSASigner(t *testing.T) {
	// Generate key pair using existing generator
	generator := &RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	// Create signer with the RSAKeyPair
	signer, err := NewRSASigner(keyPair)
	if err != nil {
		t.Fatalf("Failed to create RSA signer: %v", err)
	}

	// Sign data
	data := []byte("test data to sign")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if len(signature) == 0 {
		t.Error("Signature should not be empty")
	}

	// Verify signature using the keyPair
	err = VerifyRSA(keyPair, data, signature)
	if err != nil {
		t.Errorf("Failed to verify signature: %v", err)
	}

	// Verify with wrong data should fail
	wrongData := []byte("wrong data")
	err = VerifyRSA(keyPair, wrongData, signature)
	if err == nil {
		t.Error("Expected verification to fail with wrong data")
	}
}

// TestECDSASigner tests ECDSA signing and verification
func TestECDSASigner(t *testing.T) {
	// Generate key pair using existing generator
	generator := &ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	// Create signer with the ECCKeyPair
	signer, err := NewECDSASigner(keyPair)
	if err != nil {
		t.Fatalf("Failed to create ECDSA signer: %v", err)
	}

	// Sign data
	data := []byte("test data to sign")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if len(signature) == 0 {
		t.Error("Signature should not be empty")
	}

	// Verify signature using the keyPair
	err = VerifyECDSA(keyPair, data, signature)
	if err != nil {
		t.Errorf("Failed to verify signature: %v", err)
	}

	// Verify with wrong data should fail
	wrongData := []byte("wrong data")
	err = VerifyECDSA(keyPair, wrongData, signature)
	if err == nil {
		t.Error("Expected verification to fail with wrong data")
	}
}

// TestSignerWithNilKey tests error handling for nil keys
func TestSignerWithNilKey(t *testing.T) {
	// Test RSA signer with nil key
	_, err := NewRSASigner(nil)
	if err == nil {
		t.Error("Expected error when creating RSA signer with nil key")
	}

	// Test ECDSA signer with nil key
	_, err = NewECDSASigner(nil)
	if err == nil {
		t.Error("Expected error when creating ECDSA signer with nil key")
	}
}

// TestKeyGeneration tests key generation
func TestKeyGeneration(t *testing.T) {
	// Test RSA generation
	rsaGen := &RSAGenerator{}
	rsaKeyPair, err := rsaGen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	if rsaKeyPair.Private == nil || rsaKeyPair.Public == nil {
		t.Error("RSA key pair should not be nil")
	}

	// Test ECC generation
	eccGen := &ECCGenerator{}
	eccKeyPair, err := eccGen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}
	if eccKeyPair.Private == nil || eccKeyPair.Public == nil {
		t.Error("ECC key pair should not be nil")
	}
}

// TestRSAMarshaling tests RSA key marshaling and unmarshaling
func TestRSAMarshaling(t *testing.T) {
	generator := &RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	marshaler := NewRSAMarshaler()

	// Marshal
	publicBytes, privateBytes, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal RSA key pair: %v", err)
	}

	if len(publicBytes) == 0 || len(privateBytes) == 0 {
		t.Error("Marshaled keys should not be empty")
	}

	// Unmarshal
	unmarshaledKeyPair, err := marshaler.Unmarshal(privateBytes)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSA key pair: %v", err)
	}

	if unmarshaledKeyPair.Private == nil || unmarshaledKeyPair.Public == nil {
		t.Error("Unmarshaled key pair should not be nil")
	}

	// Test that the unmarshaled key can sign
	signer, err := NewRSASigner(unmarshaledKeyPair)
	if err != nil {
		t.Fatalf("Failed to create signer with unmarshaled key: %v", err)
	}

	data := []byte("test")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign with unmarshaled key: %v", err)
	}

	err = VerifyRSA(unmarshaledKeyPair, data, signature)
	if err != nil {
		t.Errorf("Failed to verify signature from unmarshaled key: %v", err)
	}
}

// TestECCMarshaling tests ECC key marshaling and unmarshaling
func TestECCMarshaling(t *testing.T) {
	generator := &ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	marshaler := NewECCMarshaler()

	// Encode
	publicBytes, privateBytes, err := marshaler.Encode(*keyPair)
	if err != nil {
		t.Fatalf("Failed to encode ECC key pair: %v", err)
	}

	if len(publicBytes) == 0 || len(privateBytes) == 0 {
		t.Error("Encoded keys should not be empty")
	}

	// Decode
	decodedKeyPair, err := marshaler.Decode(privateBytes)
	if err != nil {
		t.Fatalf("Failed to decode ECC key pair: %v", err)
	}

	if decodedKeyPair.Private == nil || decodedKeyPair.Public == nil {
		t.Error("Decoded key pair should not be nil")
	}

	// Test that the decoded key can sign
	signer, err := NewECDSASigner(decodedKeyPair)
	if err != nil {
		t.Fatalf("Failed to create signer with decoded key: %v", err)
	}

	data := []byte("test")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign with decoded key: %v", err)
	}

	err = VerifyECDSA(decodedKeyPair, data, signature)
	if err != nil {
		t.Errorf("Failed to verify signature from decoded key: %v", err)
	}
}

// TestSignerInterface tests that both signers implement the Signer interface
func TestSignerInterface(t *testing.T) {
	var _ Signer = (*RSASigner)(nil)
	var _ Signer = (*ECDSASigner)(nil)
}

// BenchmarkRSASigning benchmarks RSA signing performance
func BenchmarkRSASigning(b *testing.B) {
	generator := &RSAGenerator{}
	keyPair, _ := generator.Generate()
	signer, _ := NewRSASigner(keyPair)
	data := []byte("benchmark data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = signer.Sign(data)
	}
}

// BenchmarkECDSASigning benchmarks ECDSA signing performance
func BenchmarkECDSASigning(b *testing.B) {
	generator := &ECCGenerator{}
	keyPair, _ := generator.Generate()
	signer, _ := NewECDSASigner(keyPair)
	data := []byte("benchmark data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = signer.Sign(data)
	}
}
