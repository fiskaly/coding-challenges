package domain_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/util"
	"github.com/google/uuid"
)

func TestRSASigning(t *testing.T) {
	id := uuid.NewString()
	label := "Test device"
	algorithm := domain.AlgorithmRSA
	signer, err := crypto.NewRSASigner()
	failOnError(t, err)

	device := domain.NewSignatureDevice(id, label, algorithm, signer)

	ConcurrentSign(t, device, 1000)
}

func TestECCSigning(t *testing.T) {
	id := uuid.NewString()
	label := "Test device"
	algorithm := domain.AlgorithmECC
	signer, err := crypto.NewECCSigner()
	failOnError(t, err)

	device := domain.NewSignatureDevice(id, label, algorithm, signer)

	ConcurrentSign(t, device, 1000)
}

func ConcurrentSign(t *testing.T, device *domain.SignatureDevice, runCount int) {
	wg := sync.WaitGroup{}
	wg.Add(runCount)
	resultsChan := make(chan *domain.SignDataResult, runCount)
	results := make([]*domain.SignDataResult, runCount)
	data := []byte("LoremIpsum")

	for i := 0; i < runCount; i++ {
		go func() {
			r, err := device.Sign(data)
			failOnError(t, err)
			resultsChan <- r
			wg.Done()
		}()
	}

	wg.Wait()
	close(resultsChan)

	for i := 0; i < runCount; i++ {
		r := <-resultsChan
		results[i] = r
		if i == 0 {
			expectedSignedData := fmt.Sprintf("%d_%s_%s", i, data, util.EncodeToBase64String([]byte(device.Id)))
			if expectedSignedData != string(r.SignedData) {
				t.Fatalf("expected: %s | actual: %s", expectedSignedData, r.SignedData)
			}
			continue
		}

		expectedSignedData := fmt.Sprintf("%d_%s_%s", i, data, util.EncodeToBase64String(results[i-1].Signature))
		if expectedSignedData != string(r.SignedData) {
			t.Fatalf("expected: %s | actual: %s", expectedSignedData, r.SignedData)
		}
	}
}

func failOnError(t *testing.T, errs ...error) {
	err := errors.Join(errs...)
	if err != nil {
		t.Fatal(err)
	}
}
