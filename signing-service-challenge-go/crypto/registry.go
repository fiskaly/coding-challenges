package crypto

import "fmt"

func init() {
	registerToolkit("RSA", &RSA{})
	registerToolkit("ECC", &ECDSA{})
}

var toolkitRegistry = make(map[string]CryptoToolkit)

func registerToolkit(algorithm string, ct CryptoToolkit) {
	toolkitRegistry[algorithm] = ct
}

func GetToolkit(algorithm string) (CryptoToolkit, error) {
	kit, exists := toolkitRegistry[algorithm]
	if !exists {
		return nil, fmt.Errorf("no toolkit registered for %s", algorithm)
	}
	return kit, nil
}
