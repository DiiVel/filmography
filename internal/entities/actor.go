package entities

import "time"

type ActorEntity struct {
	ID       string
	Name     string
	Gender   string
	Birthday time.Time
}
