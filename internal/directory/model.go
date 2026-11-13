package directory

type Doctor struct {
	UserID       string  `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Specialty    *string `json:"specialty"`
	FacilityID   string  `json:"facility_id"`
	FacilityName string  `json:"facility_name"`
	State        string  `json:"state"`
	FeeNaira     int     `json:"fee_naira"`
	SlotMinutes  int     `json:"slot_minutes"`
	Bio          *string `json:"bio"`
}

type Search struct {
	Specialty  string
	State      string
	FacilityID string
	Name       string
}

type ProfileRequest struct {
	FeeNaira    int    `json:"fee_naira" binding:"min=0,max=10000000"`
	SlotMinutes int    `json:"slot_minutes" binding:"required,oneof=15 20 30 45 60"`
	Bio         string `json:"bio" binding:"max=500"`
	Accepting   bool   `json:"accepting"`
}

type Hours struct {
	Weekday int    `json:"weekday" binding:"min=0,max=6"`
	Starts  string `json:"starts" binding:"required,datetime=15:04"`
	Ends    string `json:"ends" binding:"required,datetime=15:04"`
}

type HoursRequest struct {
	Hours []Hours `json:"hours" binding:"max=50,dive"`
}
