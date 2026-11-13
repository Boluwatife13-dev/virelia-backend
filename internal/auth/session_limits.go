package auth

import "time"

func idleLimit(role string) time.Duration {
	switch role {
	case RolePlatformAdmin:
		return 15 * time.Minute
	case RoleHospitalAdmin:
		return 30 * time.Minute
	case RoleDoctor, RoleFrontDesk:
		return 8 * time.Hour
	default:
		return 30 * 24 * time.Hour
	}
}
