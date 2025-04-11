package models

import "gorm.io/datatypes"

type User struct {
	Id       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Requests []Request `gorm:"foreignKey:UserID" json:"requests"` // <- relacja 1:N
}

type Request struct {
	Id       uint           `json:"id"`
	UserID   uint           `json:"userid"` // <- klucz obcy
	Name     string         `json:"name"`
	PriceMin *float32       `json:"pricemin"`
	PriceMax *float32       `json:"pricemax"`
	URLs     datatypes.JSON `json:"urls"` // ✅ JSON pole kompatybilne z PostgreSQL
}
