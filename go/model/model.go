// Package model implementiert die structs,
// die in die Datenbank abgebildet werden.
package model

import "time"

// Model contains all fields for data management
// and accountability of changes for sensible data
type Model struct {
	// standard gorrm fields
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	// additional fields to make all changes to the
	// database accountable
	UpdatedBy string
	Updater   uint
	DeletedBy string
	Deleter   uint
}
