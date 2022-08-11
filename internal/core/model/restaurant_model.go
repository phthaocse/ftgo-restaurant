package model

type MenuItem struct {
	Id    int
	Name  string
	Price float64
}

type Restaurant struct {
	Id      int
	Name    string
	Address string
	Menu    []*MenuItem
}
