package entity

import "time"

type Session struct {
	ID           int
	Username     string
	Token        string
	TokenExpDate time.Time
}
