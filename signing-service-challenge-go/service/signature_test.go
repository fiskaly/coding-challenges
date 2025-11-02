package service

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

func TestSignatureServiceCreateAndSign(t *testing.T) {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	service := NewSignatureService(repo)

	createResp, err := service.CreateSignatureDevice(context.Background(), CreateSignatureDeviceRequest{
		ID:        "device-service",
		Algorithm: "rsa",
		Label:     "Checkout 1",
	})
	if err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	if createResp.ID != "device-service" {
		t.Fatalf("expected id device-service, got %s", createResp.ID)
	}

	if len(strings.TrimSpace(string(createResp.PublicKeyPEM))) == 0 {
		t.Fatal("expected public key to be returned")
	}

	signResp, err := service.SignTransaction(context.Background(), SignTransactionRequest{
		DeviceID: createResp.ID,
		Data:     "total=123",
	})
	if err != nil {
		t.Fatalf("sign transaction failed: %v", err)
	}

	if signResp.DeviceID != createResp.ID {
		t.Fatalf("expected device id %s, got %s", createResp.ID, signResp.DeviceID)
	}

	expectedPrefix := "0_total=123_"
	if !strings.HasPrefix(signResp.SignedData, expectedPrefix) {
		t.Fatalf("expected signed data to start with %q, got %q", expectedPrefix, signResp.SignedData)
	}

	expectedLast := base64.StdEncoding.EncodeToString([]byte(createResp.ID))
	if !strings.HasSuffix(signResp.SignedData, expectedLast) {
		t.Fatalf("expected signed data to end with %q, got %q", expectedLast, signResp.SignedData)
	}

	if createResp.LastSignature != expectedLast {
		t.Fatalf("expected create response last signature %q, got %q", expectedLast, createResp.LastSignature)
	}

	if signResp.Signature == "" {
		t.Fatal("expected signature to be populated")
	}

	if signResp.LastSignature != signResp.Signature {
		t.Fatalf("expected last signature to match signature, got %q vs %q", signResp.LastSignature, signResp.Signature)
	}

	if signResp.SignatureCounter != 1 {
		t.Fatalf("expected signature counter to be 1, got %d", signResp.SignatureCounter)
	}
}

func TestSignatureServiceConcurrentSigning(t *testing.T) {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	service := NewSignatureService(repo)

	if _, err := service.CreateSignatureDevice(context.Background(), CreateSignatureDeviceRequest{
		ID:        "device-concurrency",
		Algorithm: "ecc",
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	const totalSigns = 5
	var (
		wg   sync.WaitGroup
		errs = make(chan error, totalSigns)
	)

	for i := 0; i < totalSigns; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := service.SignTransaction(context.Background(), SignTransactionRequest{
				DeviceID: "device-concurrency",
				Data:     "payload-" + string(rune('a'+idx)),
			})
			if err != nil {
				errs <- err
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("sign transaction failed: %v", err)
	}

	stored, err := repo.Get(context.Background(), "device-concurrency")
	if err != nil {
		t.Fatalf("get device failed: %v", err)
	}

	if stored.SignatureCounter() != totalSigns {
		t.Fatalf("expected counter %d, got %d", totalSigns, stored.SignatureCounter())
	}
}

func TestSignatureServiceConcurrentSigningMonotonicity(t *testing.T) {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	service := NewSignatureService(repo)

	deviceID := "device-monotonic"

	if _, err := service.CreateSignatureDevice(context.Background(), CreateSignatureDeviceRequest{
		ID:        deviceID,
		Algorithm: "rsa",
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	const totalSigns = 64

	errs := make(chan error, totalSigns)
	type result struct {
		counter uint64
		data    string
		last    string
	}

	results := make(chan result, totalSigns)

	var wg sync.WaitGroup
	for i := 0; i < totalSigns; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			resp, err := service.SignTransaction(context.Background(), SignTransactionRequest{
				DeviceID: deviceID,
				Data:     "payload-" + strconv.Itoa(idx),
			})
			if err != nil {
				errs <- err
				return
			}

			results <- result{
				counter: resp.SignatureCounter,
				data:    resp.SignedData,
				last:    resp.LastSignature,
			}
		}(i)
	}

	wg.Wait()
	close(errs)
	close(results)

	for err := range errs {
		t.Fatalf("sign transaction failed: %v", err)
	}

	seenCounters := make(map[int]struct{}, totalSigns)
	lastSignatureByCounter := make(map[int]string, totalSigns) // new signature after signing
	prevSignatureByCounter := make(map[int]string, totalSigns) // previous signature embedded in payload

	for res := range results {
		prefix, payload, last := splitSignedData(res.data)
		if prefix < 0 {
			t.Fatalf("unexpected signed data format: %q", res.data)
		}

		if last == "" {
			t.Fatal("expected last signature segment to be non-empty")
		}

		if int(res.counter) != prefix+1 {
			t.Fatalf("expected response counter %d to equal payload prefix+1 (%d)", res.counter, prefix+1)
		}

		if strings.TrimSpace(payload) == "" {
			t.Fatal("expected signed payload data to be non-empty")
		}

		if _, exists := seenCounters[prefix]; exists {
			t.Fatalf("duplicate payload counter detected: %d", prefix)
		}
		seenCounters[prefix] = struct{}{}
		lastSignatureByCounter[prefix] = res.last
		prevSignatureByCounter[prefix] = last
	}

	if len(seenCounters) != totalSigns {
		t.Fatalf("expected %d unique payload counters, got %d", totalSigns, len(seenCounters))
	}

	for i := 0; i < totalSigns; i++ {
		if _, ok := seenCounters[i]; !ok {
			t.Fatalf("missing payload counter %d", i)
		}
	}

	stored, err := repo.Get(context.Background(), deviceID)
	if err != nil {
		t.Fatalf("get device failed: %v", err)
	}

	if stored.SignatureCounter() != totalSigns {
		t.Fatalf("expected final counter %d, got %d", totalSigns, stored.SignatureCounter())
	}

	initialLast := base64.StdEncoding.EncodeToString([]byte(deviceID))
	if prevSignatureByCounter[0] != initialLast {
		t.Fatalf("expected first payload to reference device id base64, got %q", prevSignatureByCounter[0])
	}

	for i := 1; i < totalSigns; i++ {
		prev, ok := prevSignatureByCounter[i]
		if !ok {
			t.Fatalf("missing previous signature link for counter %d", i)
		}

		expectedPrevious, ok := lastSignatureByCounter[i-1]
		if !ok {
			t.Fatalf("missing signature for counter %d", i-1)
		}

		if prev != expectedPrevious {
			t.Fatalf("expected counter %d payload to reference signature from counter %d", i, i-1)
		}
	}

	if stored.LastSignature() != lastSignatureByCounter[totalSigns-1] {
		t.Fatalf("expected stored last signature to equal final signature in sequence")
	}
}

func TestSignatureServiceCreateValidation(t *testing.T) {
	cases := []struct {
		name    string
		req     CreateSignatureDeviceRequest
		setup   func(*SignatureService) error
		wantErr error
	}{
		{
			name: "empty device id",
			req: CreateSignatureDeviceRequest{
				Algorithm: "rsa",
			},
			wantErr: ErrInvalidDeviceID,
		},
		{
			name: "unsupported algorithm",
			req: CreateSignatureDeviceRequest{
				ID:        "device-unsupported",
				Algorithm: "dsa",
			},
			wantErr: domain.ErrUnsupportedAlgorithm,
		},
		{
			name: "duplicate device id",
			req: CreateSignatureDeviceRequest{
				ID:        "device-dup",
				Algorithm: "rsa",
			},
			setup: func(svc *SignatureService) error {
				_, err := svc.CreateSignatureDevice(context.Background(), CreateSignatureDeviceRequest{
					ID:        "device-dup",
					Algorithm: "rsa",
				})
				return err
			},
			wantErr: domain.ErrDeviceAlreadyExists,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := persistence.NewInMemorySignatureDeviceRepository()
			svc := NewSignatureService(repo)

			if tc.setup != nil {
				if err := tc.setup(svc); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			_, err := svc.CreateSignatureDevice(context.Background(), tc.req)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func splitSignedData(signedData string) (counter int, payload string, lastSignature string) {
	lastSep := strings.LastIndex(signedData, "_")
	if lastSep == -1 {
		return -1, "", ""
	}

	lastSignature = signedData[lastSep+1:]
	prefixAndPayload := signedData[:lastSep]

	payloadSep := strings.Index(prefixAndPayload, "_")
	if payloadSep == -1 {
		return -1, "", ""
	}

	counterStr := prefixAndPayload[:payloadSep]
	payload = prefixAndPayload[payloadSep+1:]

	value, err := strconv.Atoi(counterStr)
	if err != nil {
		return -1, "", ""
	}

	return value, payload, lastSignature
}

func TestSignatureServiceSignValidation(t *testing.T) {
	cases := []struct {
		name    string
		req     SignTransactionRequest
		setup   func(*SignatureService) error
		wantErr error
	}{
		{
			name: "empty device id",
			req: SignTransactionRequest{
				DeviceID: " ",
				Data:     "payload",
			},
			wantErr: ErrInvalidDeviceID,
		},
		{
			name: "empty payload",
			req: SignTransactionRequest{
				DeviceID: "device-valid",
				Data:     "   ",
			},
			setup: func(svc *SignatureService) error {
				_, err := svc.CreateSignatureDevice(context.Background(), CreateSignatureDeviceRequest{
					ID:        "device-valid",
					Algorithm: "rsa",
				})
				return err
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "unknown device",
			req: SignTransactionRequest{
				DeviceID: "device-missing",
				Data:     "payload",
			},
			wantErr: domain.ErrDeviceNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := persistence.NewInMemorySignatureDeviceRepository()
			svc := NewSignatureService(repo)

			if tc.setup != nil {
				if err := tc.setup(svc); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			_, err := svc.SignTransaction(context.Background(), tc.req)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestSignatureServiceGetAndList(t *testing.T) {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	service := NewSignatureService(repo)

	deviceIDs := []string{"device-1", "device-2"}
	for _, id := range deviceIDs {
		if _, err := service.CreateSignatureDevice(context.Background(), CreateSignatureDeviceRequest{
			ID:        id,
			Algorithm: "rsa",
			Label:     id + "-label",
		}); err != nil {
			t.Fatalf("create device %s failed: %v", id, err)
		}
	}

	got, err := service.GetSignatureDevice(context.Background(), "device-1")
	if err != nil {
		t.Fatalf("get device failed: %v", err)
	}

	if got.ID != "device-1" || got.Label != "device-1-label" {
		t.Fatalf("unexpected device dto: %+v", got)
	}

	list, err := service.ListSignatureDevices(context.Background())
	if err != nil {
		t.Fatalf("list devices failed: %v", err)
	}

	if len(list) != len(deviceIDs) {
		t.Fatalf("expected %d devices, got %d", len(deviceIDs), len(list))
	}
}

func TestSignatureServiceMissingSigner(t *testing.T) {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	service := NewSignatureService(repo)

	device, err := domain.NewSignatureDevice("device-missing-signer", domain.AlgorithmRSA, "")
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	if err := repo.Save(context.Background(), device); err != nil {
		t.Fatalf("failed to save device: %v", err)
	}

	_, err = service.SignTransaction(context.Background(), SignTransactionRequest{
		DeviceID: "device-missing-signer",
		Data:     "payload",
	})
	if !errors.Is(err, ErrMissingSigner) {
		t.Fatalf("expected ErrMissingSigner, got %v", err)
	}
}
