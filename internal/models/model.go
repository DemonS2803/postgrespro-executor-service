package models

import "time"

type Command struct {
	Id          int       `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CompletedCommandRequest struct {
	Id          int       `json:"id"`
	Command     Command   `json:"command,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Result      string    `json:"result,omitempty"`
	Status      string    `json:"status"`
}

type CreateCommandRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
