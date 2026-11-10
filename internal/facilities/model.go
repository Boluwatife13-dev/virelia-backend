package facilities

import "time"

type Facility struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	LicenceNumber    string    `json:"licence_number"`
	LicenceExpiresAt time.Time `json:"licence_expires_at"`
	State            string    `json:"state"`
	Address          string    `json:"address"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	AdminFirstName   string    `json:"admin_first_name"`
	AdminLastName    string    `json:"admin_last_name"`
	AdminEmail       string    `json:"admin_email"`
	Status           string    `json:"status"`
	RejectionReason  *string   `json:"rejection_reason"`
	CheckedBy        *string   `json:"checked_by"`
	CheckNote        *string   `json:"check_note"`
	CreatedAt        time.Time `json:"created_at"`
}

type ApplyRequest struct {
	Name             string `json:"name" binding:"required,max=200"`
	Type             string `json:"type" binding:"required,oneof=hospital clinic primary_health_centre"`
	LicenceNumber    string `json:"licence_number" binding:"required,max=50"`
	LicenceExpiresAt string `json:"licence_expires_at" binding:"required,datetime=2006-01-02"`
	State            string `json:"state" binding:"required,max=50"`
	Address          string `json:"address" binding:"required,max=255"`
	Email            string `json:"email" binding:"required,email"`
	Phone            string `json:"phone" binding:"required,e164"`
	AdminFirstName   string `json:"admin_first_name" binding:"required,max=100"`
	AdminLastName    string `json:"admin_last_name" binding:"required,max=100"`
	AdminEmail       string `json:"admin_email" binding:"required,email"`
}

type CheckRequest struct {
	Note string `json:"note" binding:"required,max=500"`
}

type RejectRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
}
