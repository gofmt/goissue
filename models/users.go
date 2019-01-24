package models

import (
	"time"
)

type User struct {
	Id       int
	Account  string `sql:",unique"`
	NickName string
	Password string
	Type     int
	Level    int
	Status   int
	CreateAt time.Time
	UpdateAt time.Time
}
