package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type SignatureDeviceDetails struct {
	ID        string `json:"id,omitempty"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

type CreateSignatureDeviceRequest struct {
	SignatureDeviceDetails
}

type GetSignatureDeviceResponse struct {
	SignatureDeviceDetails
}

type ListSignatureDeviceResponse struct {
	IDs []domain.ID `json:"ids"`
}

type SignTransactionRequest struct {
	Data domain.Data `json:"data"`
}

type SignatureResponse struct {
	Signature  string      `json:"signature"`
	SignedData domain.Data `json:"signed_data"`
}

type api struct {
	d domain.SignatureDomain
}

func (a *api) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		panic("Create should only be called with POST")
	}
	var req CreateSignatureDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := a.d.CreateSignatureDevice(domain.Device{
		ID:        domain.ID(req.ID),
		Algorithm: req.Algorithm,
		Label:     req.Label,
	})
	if err == domain.ErrAlreadyExists {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err == domain.ErrUnsupportedAlgorithm {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *api) SignTransaction(w http.ResponseWriter, r *http.Request) {
	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := a.d.SignTransaction(domain.ID(r.PathValue("id")), req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SignatureResponse{
		Signature:  s.Signature,
		SignedData: domain.Data(s.SignedData),
	}
	json.NewEncoder(w).Encode(resp)
}

func (a *api) ListSignatureDevices(w http.ResponseWriter, r *http.Request) {
	ids, err := a.d.ListSignatureDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ListSignatureDeviceResponse{
		IDs: ids,
	})
}

func (a *api) GetSignatureDeviceDetails(w http.ResponseWriter, r *http.Request) {
	d, err := a.d.GetSignatureDeviceDetails(domain.ID(r.PathValue("id")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(GetSignatureDeviceResponse{
		SignatureDeviceDetails: SignatureDeviceDetails{
			Algorithm: d.Algorithm,
			Label:     d.Label,
		},
	})
}

func SignatureDeviceRoutes(d *domain.SignatureDomain) *http.ServeMux {
	mux := http.NewServeMux()
	a := api{
		d: *d,
	}
	mux.Handle("POST /signature-devices", http.HandlerFunc(a.CreateSignatureDevice))
	mux.Handle("GET /signature-devices", http.HandlerFunc(a.ListSignatureDevices))
	mux.Handle("GET /signature-devices/{id}", http.HandlerFunc(a.GetSignatureDeviceDetails))
	mux.Handle("POST /signature-devices/{id}/sign-transaction", http.HandlerFunc(a.SignTransaction))
	return mux
}
