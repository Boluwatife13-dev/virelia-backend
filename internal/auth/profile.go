package auth

import "context"

func (s *Service) Me(
	ctx context.Context,
	userID string,
) (*User, *Patient, error) {
	user, err := s.repo.FindUser(ctx, "id", userID)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	patient, err := s.repo.FindPatient(ctx, userID)

	return user, patient, err
}
