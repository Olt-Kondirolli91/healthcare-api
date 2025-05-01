package models

import "gorm.io/gorm"

type Patient struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FirstName string         `json:"first_name" binding:"required"`
	LastName  string         `json:"last_name"  binding:"required"`
	Email     string         `json:"email"      gorm:"unique"`
	Phone     string         `json:"phone"`
	// soft delete support
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	// relation
	Appointments []Appointment `json:"appointments,omitempty"`
}
