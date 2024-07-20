package models

import (
	"time"
)

// Project struct to describe Project object.
type Project struct {
	ID   int `db:"id" json:"id"`
	Name string `db:"name" json:"name" validate:"required"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}