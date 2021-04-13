package models

import "time"

type Model interface {
	GetID() string
	Update(deltaTime time.Duration)
	IsDeleted() bool
}
