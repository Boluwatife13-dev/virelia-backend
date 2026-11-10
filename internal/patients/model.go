package patients

import "time"

// Patient is what the front desk sees.
// No medical records, only who the person is.
type Patient struct {
	ID          string    `json:"id"`
	NetworkID   string    `json:"network_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Sex         string    `json:"sex"`
	Phone       *string   `json:"phone"`
	HasLogin    bool      `json:"has_login"`
}

// RegisterRequest is the walk-in form the
// front desk fills in. NIN is optional here.
type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required,max=100"`
	LastName    string `json:"last_name" binding:"required,max=100"`
	DateOfBirth string `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
	Sex         string `json:"sex" binding:"required,oneof=male female"`
	Phone       string `json:"phone" binding:"omitempty,e164"`
	NIN         string `json:"nin" binding:"omitempty,len=11,numeric"`

	State          string `json:"state" binding:"max=50"`
	Address        string `json:"address" binding:"max=255"`
	EmergencyName  string `json:"emergency_contact_name" binding:"max=100"`
	EmergencyPhone string `json:"emergency_contact_phone" binding:"omitempty,e164"`
}

// Search looks a patient up by one of these.
type Search struct {
	Phone     string
	NIN       string
	NetworkID string
}
