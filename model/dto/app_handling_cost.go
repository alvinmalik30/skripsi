package dto

type AppHandlingCostReq struct {
	Name    string  `json:"name"`
	Nominal float64 `json:"nominal"`
	Unit    string  `json:"unit"`
}
