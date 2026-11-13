package records

import "time"

type Entry struct {
	ID               string         `json:"id"`
	PatientID        string         `json:"patient_id"`
	FacilityID       string         `json:"facility_id"`
	FacilityName     string         `json:"facility_name"`
	AuthorID         string         `json:"author_id"`
	AuthorName       string         `json:"author_name"`
	AppointmentID    *string        `json:"appointment_id"`
	Kind             string         `json:"kind"`
	Data             map[string]any `json:"data"`
	Version          int            `json:"version"`
	ReplacesID       *string        `json:"replaces_id"`
	CorrectionReason *string        `json:"correction_reason"`
	CreatedAt        time.Time      `json:"created_at"`
}

type AddRequest struct {
	PatientNetworkID string         `json:"patient_network_id" binding:"required"`
	AppointmentID    string         `json:"appointment_id" binding:"omitempty,uuid"`
	Kind             string         `json:"kind" binding:"required,oneof=visit vitals diagnosis note"`
	Data             map[string]any `json:"data" binding:"required"`
}

type CorrectRequest struct {
	Data   map[string]any `json:"data" binding:"required"`
	Reason string         `json:"reason" binding:"required,max=300"`
}
