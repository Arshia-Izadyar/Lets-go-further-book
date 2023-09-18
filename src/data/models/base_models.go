package models

import "time"

type Movie struct {
	ID        int
	CreatedAt time.Time
	Title     string
	Year      int32
	Genres    []string
	Runtime   Runtime
	Version   int32
}
