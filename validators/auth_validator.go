package validators

type Register struct {
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,e164"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SendOTP struct {
	Channel string `json:"channel" binding:"required,oneof=email phone"`
}

type VerifyOTP struct {
	Channel string `json:"channel" binding:"required,oneof=email phone"`
	Code    string `json:"code" binding:"required,len=6,numeric"`
}

type ForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPassword struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6,numeric"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type ChangePassword struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
}