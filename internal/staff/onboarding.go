package staff

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SendInvite(ctx context.Context, in Invite) error {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	in.Email = strings.ToLower(in.Email)

	err := s.repo.CreateInvite(ctx, in, hash(token))
	if err != nil {
		return err
	}

	link := s.appURL + "/accept-invite?token=" + token

	body := `<p>You've been invited to Virelia.</p>` +
		`<p><a href="` + link + `">Accept your invite</a></p>` +
		`<p>This link expires in 72 hours.</p>`

	return s.mailer.Send(ctx, in.Email, "Your Virelia invite", body)
}

func (s *Service) Accept(ctx context.Context, req AcceptRequest) error {
	pw, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)

	in, err := s.repo.Accept(ctx, hash(req.Token), req.Phone, string(pw))
	if err != nil {
		return err
	}

	admins, _ := s.repo.AdminEmails(ctx, in.FacilityID)

	name := html.EscapeString(in.FirstName + " " + in.LastName)

	body := "<p>" + name + " has joined your hospital.</p>"

	s.mailer.SendAll(ctx, admins, "New staff member joined", body)

	return nil
}
