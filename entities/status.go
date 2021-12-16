package entities

type EntityState uint8

const (
	StateOff EntityState = iota
	StateOn
)