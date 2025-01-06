package models

type Product struct {
	Id    string  `json:"Id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
