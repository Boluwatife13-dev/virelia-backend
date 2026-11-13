package staff

type Member struct {
	UserID    string  `json:"user_id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Role      string  `json:"role"`
	Status    string  `json:"status"`
	Licence   *string `json:"licence_number"`
	Specialty *string `json:"specialty"`
}

type Invite struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	Role      string `json:"role" binding:"required,oneof=doctor front_desk hospital_admin"`
	Licence   string `json:"licence_number" binding:"required_if=Role doctor,max=50"`
	Specialty string `json:"specialty" binding:"max=100"`

	FacilityID string `json:"-"`
	InvitedBy  string `json:"-"`
}

type AcceptRequest struct {
	Token    string `json:"token" binding:"required"`
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
