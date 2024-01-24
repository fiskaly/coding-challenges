package crypto

import "testing"

func TestECCAlgorithm_Name(t *testing.T) {
	ecc := ECCAlgorithm{}
	name := ecc.Name()

	if name != "ECC" {
		t.Errorf("expected ECC, got: %s", name)
	}
}

func TestECCAlgorithm_GenerateEncodedPrivateKey(t *testing.T) {
	ecc := ECCAlgorithm{}
	encodedPrivateKey, err := ecc.GenerateEncodedPrivateKey()

	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	_, err = NewECCMarshaler().Decode(encodedPrivateKey)
	if err != nil {
		t.Errorf("decode of generated private key failed: %s", err)
	}
}
