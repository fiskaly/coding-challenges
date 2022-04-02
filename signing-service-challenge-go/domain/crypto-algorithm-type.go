package domain

type CryptoAlgorithmType string

const (
	RSA CryptoAlgorithmType = "RSA"
	ECC                     = "ECC"
)

func (algorithmType CryptoAlgorithmType) String() string {
	switch algorithmType {
	case RSA:
		return "RSA"
	case ECC:
		return "ECC"
	}
	return "unknown"
}
