package api

type DeviceListResponse struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	Algorithm string `json:"algorithm"`
}

type DeviceByIdResponse struct {
	Id               string `json:"id"`
	Label            string `json:"label"`
	Algorithm        string `json:"algorithm"`
	SignatureCounter int    `json:"signature_counter"`
}
