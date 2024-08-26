package model

import (
	"time"
)

type UserAccount struct {
	Id             string    `json:"id"`
	Username       string    `json:"userName"`
	Password       string    `json:"password"`
	Email          string    `json:"email"`
	StartDate      time.Time `json:"start_date"`
	CaducityDate   time.Time `json:"expired_at"` // Fecha de expiración
	Status         int       `json:"status"`
	ChangePassword bool      `json:"changePassword"`
}
