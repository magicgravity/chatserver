package pojo

import "time"

type User struct{
	Id 				string	`json:"id"`
	UserName		string	`json:"UserName"`
	UserId			string	`json:"UserId"`
	Email			string	`json:"Email"`
	MobileNo		string	`json:"MobileNo"`
	Address			string	`json:"Address"`
	Sex				int		`json:"Sex"`
	Introduce		string	`json:"Introduce"`
	Avatar			string	`json:"Avatar"`
	BgImgUrl		string	`json:"BgImgUrl"`
	Job				string	`json:"Job"`
	City			string	`json:"City"`
	Country			string	`json:"Country"`
	Password 		string 	`json:"Password"`
	CreateTime		time.Time
	UpdateTime		time.Time
}