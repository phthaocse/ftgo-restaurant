package model

type MenuItem struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Restaurant struct {
	Name    string      `json:"name"`
	Address string      `json:"address"`
	Menu    []*MenuItem `json:"menu"`
}
