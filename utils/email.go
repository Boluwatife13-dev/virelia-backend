package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const brevoURL = "https://api.brevo.com/v3/smtp/email"

type Mailer struct {
	APIKey string
	From   string
}

func (m Mailer) Send(
	ctx context.Context, to, subject, html string,
) error {
	body, err := json.Marshal(map[string]any{
		"sender":      map[string]string{"name": "Virelia", "email": m.From},
		"to":          []map[string]string{{"email": to}},
		"subject":     subject,
		"htmlContent": html,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, brevoURL, bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("api-key", m.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("brevo returned %d", res.StatusCode)
	}

	return nil
}