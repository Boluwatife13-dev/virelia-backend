package audit

import "time"

type Entry struct {
	ActorID    string
	ActorRole  string
	FacilityID string
	Action     string
	PatientID  string
	TargetID   string
	Details    map[string]any
	IP         string
	UserAgent  string
}

type Record struct {
	ID         int64          `json:"id"`
	ActorID    *string        `json:"actor_id"`
	ActorRole  *string        `json:"actor_role"`
	FacilityID *string        `json:"facility_id"`
	Action     string         `json:"action"`
	PatientID  *string        `json:"patient_id"`
	TargetID   *string        `json:"target_id"`
	Details    map[string]any `json:"details"`
	IP         *string        `json:"ip"`
	CreatedAt  time.Time      `json:"created_at"`
}
