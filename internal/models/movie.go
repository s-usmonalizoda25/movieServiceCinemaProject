package models

import "time"

type Movie struct {
	ID          int64
	Title       string
	Description string
	Duration    int32
	AgeLimit    int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
