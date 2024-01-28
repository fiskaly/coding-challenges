package crypto

import "testing"

func TestFindKeyPairGenerator(t *testing.T) {
	t.Run("returns found: false when generator does not exist", func(t *testing.T) {
		_, found := FindKeyPairGenerator("INVALID")
		if found {
			t.Error("expected found: false")
		}
	})

	t.Run("returns generator when it exists", func(t *testing.T) {
		algorithmName := "RSA"
		generator, found := FindKeyPairGenerator(algorithmName)
		if !found {
			t.Error("expected found: true")
		}

		if generator.AlgorithmName() != algorithmName {
			t.Errorf("expected %s, got %s", algorithmName, generator.AlgorithmName())
		}
	})
}
