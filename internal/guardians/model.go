package guardians

import "time"

type Child struct {
	FirstName    string `json:"first_name" binding:"required,max=100"`
	LastName     string `json:"last_name" binding:"required,max=100"`
	DateOfBirth  string `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
	Sex          string `json:"sex" binding:"required,oneof=male female"`
	NIN          string `json:"nin" binding:"omitempty,len=11,numeric"`
	Relationship string `json:"relationship" binding:"required,oneof=mother father guardian other"`
}

type FrontDeskChild struct {
	GuardianNetworkID string `json:"guardian_network_id" binding:"required"`
	Child
}

type Link struct {
	ID           string    `json:"id"`
	NetworkID    string    `json:"network_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Relationship string    `json:"relationship"`
	Verified     bool      `json:"verified"`
	ExpiresAt    time.Time `json:"expires_at"`
}
