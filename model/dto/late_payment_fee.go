package dto

type LatePaymentFeeReq struct {
	Name    string  `json:"name"`
	Nominal float64 `json:"nominal"`
	Unit    string  `json:"unit"`
}
