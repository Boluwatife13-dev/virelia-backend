package auth

import "context"

// These let the keys module log platform admins in
// without reaching into auth's private code.

// ParseChallenge reads a step-1 login token.
func (s *Service) ParseChallenge(token string) (string, string, error) {
	claims, err := parseToken(token, challengeIssuer, s.secret)
	if err != nil {
		return "", "", err
	}

	return claims.Subject, claims.Role, nil
}

// IssueFor gives a fully logged-in token pair.
func (s *Service) IssueFor(
	ctx context.Context,
	userID, agent, ip string,
) (any, error) {
	return s.issueTokens(ctx, userID, "", agent, ip)
}

// HasKey says if a user has a security key.
func (r *Repository) HasKey(ctx context.Context, userID string) bool {
	var has bool

	r.db.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM security_keys WHERE user_id = $1
        )`, userID).Scan(&has)

	return has
}
