package keys

import "github.com/go-webauthn/webauthn/webauthn"

// admin is a platform admin as the webauthn
// library needs to see them.
type admin struct {
	id    string
	email string
	keys  []webauthn.Credential
}

func (a *admin) WebAuthnID() []byte {
	return []byte(a.id)
}

func (a *admin) WebAuthnName() string {
	return a.email
}

func (a *admin) WebAuthnDisplayName() string {
	return a.email
}

func (a *admin) WebAuthnCredentials() []webauthn.Credential {
	return a.keys
}
