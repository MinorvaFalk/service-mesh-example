package model

import "encoding/json"

type HTTPResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (h *HTTPResponse) Byte() []byte {
	b, _ := json.Marshal(h)
	return b
}
