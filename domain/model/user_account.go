package model

import (
	"time"
)

type UserAccount struct {
	Id             string    `json:"id"`
	Username       string    `json:"userName"`
	Password       string    `json:"password"`
	Email          string    `json:"email"`
	StartDate      time.Time `json:"startDate"`
	CaducityDate   time.Time `json:"caducity_Date"`
	Status         int       `json:"status"`
	ChangePassword bool      `json:"changePassword"`
}
