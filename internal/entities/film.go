package entities

import "time"

// Film model
// @SWG.Model
type FilmEntity struct {
	ID          string
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float64
	Actors      []ActorEntity
}
