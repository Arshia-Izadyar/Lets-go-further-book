package dto

import (
	"clean_api/src/data/models"
	"time"
)

type CreateMovie struct {
	Title   string         `json:"title"`
	Year    int32          `json:"year"`
	Genres  []string       `json:"genres"`
	Runtime models.Runtime `json:"runtime"`
}

type UpdateMovie struct {
	Title   string         `json:"title,omitempty"`
	Year    int32          `json:"year,omitempty"`
	Genres  []string       `json:"genres,omitempty"`
	Runtime *models.Runtime `json:"runtime,omitempty"`
}

type MovieResponse struct {
	Id        int            `json:"id"`
	Title     string         `json:"title,omitempty"`
	Year      int32          `json:"year,omitempty"`
	Genres    []string       `json:"genres,omitempty"`
	Runtime   models.Runtime `json:"runtime,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	Version   int            `json:"version"`
}

