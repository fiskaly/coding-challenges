package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

func TestDeviceHandlerCreateAndSignFlow(t *testing.T) {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	svc := service.NewSignatureService(repo)
	handler := NewDeviceHandler(svc)

	createBody := bytes.NewBufferString(`{"id":"device-test","algorithm":"rsa","label":"POS"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices", createBody)
	createRec := httptest.NewRecorder()

	handler.HandleCollection(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", createRec.Code)
	}

	var createResp Response
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	signBody := bytes.NewBufferString(`{"data":"total=42"}`)
	signReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices/device-test/sign", signBody)
	signRec := httptest.NewRecorder()

	handler.HandleResource(signRec, signReq)

	if signRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", signRec.Code)
	}

	var signResp Response
	if err := json.Unmarshal(signRec.Body.Bytes(), &signResp); err != nil {
		t.Fatalf("failed to unmarshal sign response: %v", err)
	}

	if signResp.Data == nil {
		t.Fatal("expected sign response data")
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices", nil)
	listRec := httptest.NewRecorder()
	handler.HandleCollection(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listRec.Code)
	}
}

func TestMapError(t *testing.T) {
	cases := []struct {
		name        string
		err         error
		wantStatus  int
		wantMessage string
	}{
		{
			name:       "nil error",
			err:        nil,
			wantStatus: http.StatusOK,
		},
		{
			name:        "invalid device id",
			err:         service.ErrInvalidDeviceID,
			wantStatus:  http.StatusBadRequest,
			wantMessage: service.ErrInvalidDeviceID.Error(),
		},
		{
			name:        "domain validation error",
			err:         domain.ErrUnsupportedAlgorithm,
			wantStatus:  http.StatusBadRequest,
			wantMessage: domain.ErrUnsupportedAlgorithm.Error(),
		},
		{
			name:        "conflict",
			err:         domain.ErrDeviceAlreadyExists,
			wantStatus:  http.StatusConflict,
			wantMessage: domain.ErrDeviceAlreadyExists.Error(),
		},
		{
			name:        "not found",
			err:         domain.ErrDeviceNotFound,
			wantStatus:  http.StatusNotFound,
			wantMessage: domain.ErrDeviceNotFound.Error(),
		},
		{
			name:        "missing signer",
			err:         service.ErrMissingSigner,
			wantStatus:  http.StatusServiceUnavailable,
			wantMessage: service.ErrMissingSigner.Error(),
		},
		{
			name:        "internal error fallback",
			err:         errors.New("boom"),
			wantStatus:  http.StatusInternalServerError,
			wantMessage: http.StatusText(http.StatusInternalServerError),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, messages := mapError(tc.err)

			if status != tc.wantStatus {
				t.Fatalf("expected status %d, got %d", tc.wantStatus, status)
			}

			switch {
			case tc.wantMessage == "":
				if len(messages) != 0 {
					t.Fatalf("expected no messages, got %v", messages)
				}
			default:
				if len(messages) != 1 || messages[0] != tc.wantMessage {
					t.Fatalf("expected message %q, got %v", tc.wantMessage, messages)
				}
			}
		})
	}
}
