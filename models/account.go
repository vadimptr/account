package models

type Account struct {
	Name   string `gorm:"column:name;primary_key"`
	Amount int    `gorm:"column:amount;not null;default:0"`
}
