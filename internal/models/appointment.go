package models

import "time"

type Appointment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PatientID uint      `json:"patient_id" binding:"required"`
	DateTime  time.Time `json:"date_time"  binding:"required"`
	Notes     string    `json:"notes"`
}
