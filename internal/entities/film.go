package entities

import "time"

type FilmEntity struct {
	ID          string
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float64
}
