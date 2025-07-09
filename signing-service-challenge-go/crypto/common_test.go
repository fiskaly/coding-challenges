package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToSA(t *testing.T) {

	// test not-implemented SingingAlgorithm
	_, err := stringToSA("RUSSO")
	assert.NotNil(t, err)

	// test working cases
	ecc, err := stringToSA("ECC")
	assert.Nil(t, err)
	assert.Equal(t, ecc, ECC)
	rsa, err := stringToSA("RSA")
	assert.Nil(t, err)
	assert.Equal(t, rsa, RSA)
}

func TestSAToString(t *testing.T) {

	// test not-implemented SingingAlgorithm
	var RUSSO SigningAlgorithm = 4
	str := RUSSO.String()
	assert.Equal(t, str, "")

	// test correct strings
	str = ECC.String()
	assert.Equal(t, str, "ECC")
	str = RSA.String()
	assert.Equal(t, str, "RSA")
}
