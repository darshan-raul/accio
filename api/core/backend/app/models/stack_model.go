package models

import (
	"time"
)

// Stack struct to describe Stack object.
type Stack struct {
	ID   int `db:"id" json:"id"`
	Name string `db:"name" json:"name" validate:"required"`
	ProjectId int `db:"project_id" json:"project_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

}

type CreateStackRequest struct {
	StackName   string `json:"stack_name" validate:"required,min=3,max=32"`
	ProjectName string `json:"project_name" validate:"required,min=3,max=32"`
}
