package appointments

import "time"

var wat = time.FixedZone("WAT", 60*60)

type Appointment struct {
	ID           string    `json:"id"`
	PatientID    string    `json:"patient_id"`
	PatientName  string    `json:"patient_name"`
	NetworkID    string    `json:"network_id"`
	DoctorID     string    `json:"doctor_id"`
	DoctorName   string    `json:"doctor_name"`
	FacilityID   string    `json:"facility_id"`
	FacilityName string    `json:"facility_name"`
	StartsAt     time.Time `json:"starts_at"`
	EndsAt       time.Time `json:"ends_at"`
	Status       string    `json:"status"`
	Reason       *string   `json:"reason"`
}

type Slot struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type BookRequest struct {
	DoctorID string    `json:"doctor_id" binding:"required,uuid"`
	StartsAt time.Time `json:"starts_at" binding:"required"`
	Reason   string    `json:"reason" binding:"max=300"`
}

type DeskBookRequest struct {
	PatientNetworkID string `json:"patient_network_id" binding:"required"`
	BookRequest
}

type CancelRequest struct {
	Reason string `json:"reason" binding:"max=300"`
}
