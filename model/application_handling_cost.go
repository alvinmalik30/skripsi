package model

type AppHandlingCost struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Nominal float64 `json:"nominal"`
	Unit    string  `json:"unit"`
}
