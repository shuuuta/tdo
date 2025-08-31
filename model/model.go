package model

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ProjectPath string `json:"project_path"`
	Tasks       []Task `json:"tasks"`
}
