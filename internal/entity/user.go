package entity

import "time"

type User struct {
	ID                int
	Email             string
	Username          string
	Password          string
	CreationDate      time.Time
	CreationDateFront string
}
