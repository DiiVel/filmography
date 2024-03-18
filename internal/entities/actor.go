package entities

import "time"

// Actor model
// @SWG.Model
type ActorEntity struct {
	ID       string
	Name     string
	Gender   string
	Birthday time.Time
}
